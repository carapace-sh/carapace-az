# Extensions

The Azure CLI extension system — how extensions work, are discovered, installed, and managed.

> **Source of truth**: <https://learn.microsoft.com/en-us/cli/azure/azure-cli-extensions-overview> and <https://learn.microsoft.com/en-us/cli/azure/extension>. For architecture details, see [architecture.md](architecture.md).

## Overview

Extensions are **Python wheels** (`.whl` files) that are not shipped with the core CLI but are dynamically loaded as CLI commands. They provide:

- **Preview/experimental commands** before promotion to core
- **Service-specific commands** that are too niche for core
- **Third-party integrations** (e.g., `az devops`, `az ml`, `az terraform`)

Key characteristics:
- Extensions are **Python wheels** installed from a Microsoft extension index or a custom URL/path
- Extensions **do NOT auto-update** with the CLI — they must be updated separately
- Extensions **cannot depend on each other**
- Extensions **can override** existing core commands (with a warning)
- Extension module names always start with `azext_`

## Extension Management Commands

| Command | Description |
|---------|-------------|
| `az extension list` | List installed extensions |
| `az extension list-available` | List available extensions from the index |
| `az extension add --name <name>` | Install an extension from the index |
| `az extension add --source <URL-or-path>` | Install from a specific wheel file |
| `az extension remove --name <name>` | Remove an extension |
| `az extension update --name <name>` | Update an extension |
| `az extension show --name <name>` | Show extension details |

```bash
# List installed extensions
az extension list --output table

# Install an extension
az extension add --name ml

# Install from a specific URL
az extension add --source https://example.com/azext_custom-1.0.0-py3-none-any.whl

# Update an extension
az extension update --name ml

# Remove an extension
az extension remove --name ml
```

## Dynamic Installation

Since Azure CLI v2.10.0, when you run an extension command that isn't installed, `az` can auto-install it.

| Mode | Config Value | Behavior |
|------|-------------|----------|
| Off (default) | `no` | Error: command not found |
| Prompt | `yes_prompt` | Prompts before installing |
| No prompt | `yes_without_prompt` | Installs silently |

```bash
# Enable dynamic install with prompt
az config set extension.use_dynamic_install=yes_prompt

# Enable dynamic install without prompt (for automation)
az config set extension.use_dynamic_install=yes_without_prompt

# Disable dynamic install
az config set extension.use_dynamic_install=no
```

Related setting:

```bash
# Whether to continue running the command after dynamic install (default: true)
az config set extension.run_after_dynamic_install=true
```

## Extension Discovery and Loading

### Install Locations

Extensions are discovered from multiple directories:

| Location | Environment Variable | Default Path |
|----------|---------------------|-------------|
| User extensions | `AZURE_EXTENSION_DIR` | `~/.azure/cliextensions/` (Linux/macOS), `%USERPROFILE%\.azure\cliextensions\` (Windows) |
| System extensions | (fixed) | `<site-packages>/azure-cli-extensions/` |
| Dev extensions | (code-based) | Arbitrary paths for development |

### Extension Types

| Type | Description |
|------|-------------|
| `WheelExtension` | Installed from a `.whl` file; found by scanning extension dirs for `*.dist-info` or `*.egg-info` |
| `DevExtension` | Installed in development mode; found by searching `DEV_EXTENSION_SOURCES` recursively |

### Loading Process

1. `MainCommandsLoader.load_command_table()` checks the command index
2. If a command requires an extension, the extension module is loaded
3. Extension's directory is appended to `sys.path`
4. The `azext_*` module is imported
5. Extension's `COMMAND_LOADER_CLS.load_command_table()` is called
6. Commands are merged into the command table
7. `cmd.command_source = ExtensionCommandSource(...)` tracks the source

### Always-Loaded Extensions

Some extensions are always loaded regardless of the command being run:

```python
ALWAYS_LOADED_EXTENSIONS = ['azext_ai_examples', 'azext_next']
```

## Extension Metadata

Each extension has an `azext_metadata.json` file with compatibility information:

| Field | Description |
|-------|-------------|
| `azext.minCliCoreVersion` | Minimum CLI core version required |
| `azext.maxCliCoreVersion` | Maximum CLI core version supported |
| `azext.isPreview` | Whether the extension is in preview |
| `azext.isExperimental` | Whether the extension is experimental |

Compatibility is checked during loading — if the CLI core version doesn't meet the extension's requirements, the extension is not loaded (with a warning).

## Private Extension Index

Since CLI v2.20.0, you can configure a custom extension index:

```bash
az config set extension.index_url=https://mycompany.com/extensions/
```

This allows organizations to host internal extensions that aren't published to the public Microsoft index.

## Core vs Extension

Many command groups exist in **both** core and as extensions. This is indicated by "Core and Extension" in the command reference.

| Status | Description |
|--------|-------------|
| **Core** | Shipped with the CLI, updated with the CLI |
| **Extension** | Must be installed separately, updated separately |
| **Core and Extension** | Core provides base commands; extension adds more |

When an extension overrides a core command, the CLI warns you:

```
The installed extension 'ml' is in preview and adds or overwrites 'az ml' command group. It can be removed using 'az extension remove --name ml'.
```

## Common Extensions

### Popular Microsoft-Published Extensions

| Extension | Commands | Description |
|-----------|---------|-------------|
| `ml` | `az ml` | Azure Machine Learning v2 |
| `devops` | `az devops`, `az boards`, `az repos`, `az pipelines` | Azure DevOps |
| `graph` | `az graph` | Azure Resource Graph queries |
| `ssh` | `az ssh` | SSH into VMs using Microsoft Entra ID certificates |
| `storage-preview` | `az storage` (additional commands) | Preview storage features |
| `azure-firewall` | `az network firewall` | Azure Firewall management |
| `log-analytics` | `az monitor log-analytics` (data-plane query) | Log Analytics data-plane queries |
| `interactive` | `az interactive` | Interactive shell (auto-installed) |
| `connection` | `az connection` | Service Connector |
| `ai-examples` | `az ai-examples` | AI-powered help examples |
| `next` | `az next` | Command recommendations |
| `scenario` | `az scenario` | E2E scenario guidance |
| `terraform` | `az terraform` | Azure Terraform integration |
| `k8s-extension` | `az k8s-extension` | Kubernetes extensions |
| `connectedk8s` | `az connectedk8s` | Connected Kubernetes clusters |
| `customlocation` | `az customlocation` | Custom locations |
| `databricks` | `az databricks` | Databricks workspaces |
| `sentinel` | `az sentinel` | Microsoft Sentinel |
| `datafactory` | `az datafactory` | Azure Data Factory |
| `containerapp` | `az containerapp` (additional commands) | Container Apps preview features |

### Installing Multiple Extensions

```bash
# Install several at once
az extension add --name ml
az extension add --name devops
az extension add --name graph

# List all installed
az extension list --output table
```

## Extension Development

Extensions are developed as Python packages with the `azext_` prefix:

```
azext_myextension/
├── azext_myextension/
│   ├── __init__.py        # Exports COMMAND_LOADER_CLS
│   ├── commands.py        # Command registrations
│   ├── params.py          # Argument definitions
│   ├── help.py            # YAML help strings
│   ├── custom.py          # Handler functions
│   └── azext_metadata.json  # Extension metadata
├── setup.py
└── README.md
```

The extension's `COMMAND_LOADER_CLS` must inherit from `AzCommandsLoader` and implement `load_command_table()` and `load_arguments()`, just like core command modules.

Build and install:

```bash
# Build the wheel
python setup.py bdist_wheel

# Install locally
az extension add --source ./dist/azext_myextension-1.0.0-py3-none-any.whl

# Or use development mode
pip install -e .
```

## References

- Extensions overview: <https://learn.microsoft.com/en-us/cli/azure/azure-cli-extensions-overview>
- Extension commands: <https://learn.microsoft.com/en-us/cli/azure/extension>
- Available extensions: <https://learn.microsoft.com/en-us/cli/azure/azure-cli-extensions-list>
- Extension development: <https://github.com/Azure/azure-cli/blob/dev/doc/authoring_extensions.md>

## Related Skills

- For CLI architecture and module loading, see [architecture.md](architecture.md)
- For configuration of dynamic install settings, see [config.md](config.md)
- For extension auto-update gotchas, see [gotchas.md](gotchas.md)
