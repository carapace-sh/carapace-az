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

    return {
        'name': name,
        'options': options,
        'help': settings.get('help', ''),
        'required': settings.get('required', False),
        'choices': list(settings['choices']) if settings.get('choices') else None,
        'type': str(settings.get('type', '')),
        'nargs': str(settings.get('nargs', '')) or None,
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
        'group': name.split()[0] if ' ' in name or len(name.split()) > 1 else name,
    }


def serialize_group(name, group):
    return {
        'help': getattr(group, 'group_help', '') or '',
        'groups': {},
    }


def main():
    try:
        from azure.cli.core import AzCli, MainCommandsLoader
        from azure.cli.core.commands import AzCliCommandInvoker
        from azure.cli.core.parser import AzCliCommandParser
        from azure.cli.core._config import GLOBAL_CONFIG_DIR, ENV_VAR_PREFIX
        from azure.cli.core._help import AzCliHelp
        from azure.cli.core._output import AzOutputProducer
        from azure.cli.core.azlogging import AzCliLogging
        from azure.cli.core.file_util import create_invoker_and_load_cmds_and_args
    except ImportError as e:
        print(f'Error: Azure CLI not available: {e}', file=sys.stderr)
        sys.exit(1)

    cli = AzCli(
        cli_name='az',
        config_dir=GLOBAL_CONFIG_DIR,
        config_env_var_prefix=ENV_VAR_PREFIX,
        commands_loader_cls=MainCommandsLoader,
        invocation_cls=AzCliCommandInvoker,
        parser_cls=AzCliCommandParser,
        logging_cls=AzCliLogging,
        output_cls=AzOutputProducer,
        help_cls=AzCliHelp,
    )

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

    output = {
        'cli': {
            'name': 'az',
            'version': '',
        },
        'commands': commands,
        'groups': groups,
    }

    json.dump(output, sys.stdout, indent=2, default=str)
    print()


if __name__ == '__main__':
    main()
