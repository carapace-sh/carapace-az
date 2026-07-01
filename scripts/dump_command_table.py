#!/usr/bin/env python3
"""Dump the Azure CLI command table as JSON for carapace spec generation.

Uses Knack's command_table API to extract all commands, arguments,
and command groups with rich metadata. Descriptions are sourced from
the knack help_files registry which contains YAML help text authored
in the az CLI's _help.py modules.
"""
import json
import math
import sys


def sanitize_default(obj):
    """Convert non-serializable values to strings, handling NaN/Inf."""
    if isinstance(obj, float):
        if math.isnan(obj) or math.isinf(obj):
            return None
        return obj
    return str(obj)


def parse_help_entry(entry):
    """Parse a knack help YAML string into a dict."""
    if not entry:
        return {}
    try:
        import yaml
        if isinstance(entry, str):
            return yaml.safe_load(entry) or {}
        return entry if isinstance(entry, dict) else {}
    except Exception:
        return {}


def get_help(helps, key):
    """Get parsed help for a command or group name."""
    return parse_help_entry(helps.get(key))


def get_param_help(helps, cmd_name, options):
    """Get help text for a parameter from the command's help entry."""
    cmd_help = get_help(helps, cmd_name)
    params = cmd_help.get('parameters', [])
    for param in params:
        param_name = param.get('name', '')
        for opt in options:
            if param_name and param_name in [opt, opt.lstrip('-')]:
                return param.get('short-summary', '') or param.get('long-summary', '')
    return ''


def serialize_argument(name, arg, helps, cmd_name):
    settings = arg.type.settings
    options = settings.get('options_list', [])
    if isinstance(options, str):
        options = [options]

    help_val = settings.get('help', '')
    if not isinstance(help_val, str):
        help_val = str(help_val) if help_val else ''

    if not help_val:
        help_val = get_param_help(helps, cmd_name, options)

    choices = settings.get('choices')
    if callable(choices):
        try:
            choices = choices()
        except Exception:
            choices = None
    if choices:
        choices = list(choices)

    nargs = settings.get('nargs')
    if nargs is not None:
        nargs = str(nargs)

    default = settings.get('default', None)
    if isinstance(default, float) and (math.isnan(default) or math.isinf(default)):
        default = None

    return {
        'name': name,
        'options': options,
        'help': help_val,
        'required': settings.get('required', False),
        'choices': choices,
        'type': str(settings.get('type', '')),
        'nargs': nargs,
        'default': default,
        'metavar': settings.get('metavar', None),
    }


def serialize_command(name, cmd, helps):
    description = cmd.description
    if callable(description):
        try:
            description = description()
        except Exception:
            description = ''

    if not description:
        az_help = getattr(cmd, 'AZ_HELP', None) or getattr(cmd, 'help', None)
        if isinstance(az_help, dict):
            description = az_help.get('short-summary', '') or az_help.get('long-summary', '')

    if not description:
        help_data = get_help(helps, name)
        description = help_data.get('short-summary', '') or help_data.get('long-summary', '')

    arguments = []
    for arg_name, arg in cmd.arguments.items():
        arguments.append(serialize_argument(arg_name, arg, helps, name))

    return {
        'description': description or '',
        'arguments': arguments,
        'group': name.split()[0],
    }


def serialize_group(name, group, helps):
    help_data = get_help(helps, name)
    return {
        'help': help_data.get('short-summary', '') or '',
        'groups': {},
    }


def main():
    try:
        from azure.cli.core import get_default_cli
        from azure.cli.core.file_util import create_invoker_and_load_cmds_and_args
    except ImportError as e:
        print(f'Error: Azure CLI not available: {e}', file=sys.stderr)
        sys.exit(1)

    cli = get_default_cli()

    try:
        from azure.cli.core import EVENT_FAILED_EXTENSION_LOAD
        def _extension_failed_handler(_, event_data):
            print(f'Warning: failed to load extension: {event_data}', file=sys.stderr)
        cli.register_event(EVENT_FAILED_EXTENSION_LOAD, _extension_failed_handler)
    except ImportError:
        pass

    create_invoker_and_load_cmds_and_args(cli)

    from knack.help_files import helps

    invoker = cli.invocation
    command_table = invoker.commands_loader.command_table
    command_group_table = invoker.commands_loader.command_group_table

    commands = {}
    for name, cmd in command_table.items():
        commands[name] = serialize_command(name, cmd, helps)

    groups = {}
    for name, group in command_group_table.items():
        groups[name] = serialize_group(name, group, helps)

    try:
        from azure.cli.core import __version__ as az_version
    except ImportError:
        az_version = ''

    output = {
        'cli': {
            'name': 'az',
            'version': az_version,
        },
        'commands': commands,
        'groups': groups,
    }

    json.dump(output, sys.stdout, indent=2, allow_nan=False, default=sanitize_default)
    print()


if __name__ == '__main__':
    main()
