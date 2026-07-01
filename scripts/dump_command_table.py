#!/usr/bin/env python3
"""Dump the Azure CLI command table as JSON for carapace spec generation.

Uses Knack's command_table API to extract all commands, arguments,
and command groups with rich metadata.
"""
import json
import sys


def serialize_argument(name, arg):
    settings = arg.type.settings
    options = settings.get('options_list', [])
    if isinstance(options, str):
        options = [options]

    help_val = settings.get('help', '')
    if not isinstance(help_val, str):
        help_val = str(help_val) if help_val else ''

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

    return {
        'name': name,
        'options': options,
        'help': help_val,
        'required': settings.get('required', False),
        'choices': choices,
        'type': str(settings.get('type', '')),
        'nargs': nargs,
        'default': settings.get('default', None),
        'metavar': settings.get('metavar', None),
    }


def serialize_command(name, cmd):
    description = cmd.description
    if callable(description):
        try:
            description = description()
        except Exception:
            description = ''

    arguments = []
    for arg_name, arg in cmd.arguments.items():
        arguments.append(serialize_argument(arg_name, arg))

    return {
        'description': description or '',
        'arguments': arguments,
        'group': name.split()[0],
    }


def serialize_group(name, group):
    return {
        'help': getattr(group, 'group_help', '') or '',
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

    invoker = cli.invocation
    command_table = invoker.commands_loader.command_table
    command_group_table = invoker.commands_loader.command_group_table

    commands = {}
    for name, cmd in command_table.items():
        commands[name] = serialize_command(name, cmd)

    groups = {}
    for name, group in command_group_table.items():
        groups[name] = serialize_group(name, group)

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

    json.dump(output, sys.stdout, indent=2, default=str)
    print()


if __name__ == '__main__':
    main()
