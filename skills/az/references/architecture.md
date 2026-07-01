# Azure CLI Architecture

The Azure CLI's structural design, command dispatch model, and internal framework.

## Command Pattern

```
az [command group] [subgroup] [operation] [positional args] [flags]
```

Examples:
- `az vm create --name MyVm --resource-group MyGroup --image Ubuntu2204`
- `az network vnet subnet create --name MySubnet --resource-group MyGroup --vnet-name MyVnet`
- `az storage account list --output table`

## Hierarchical Command Tree

```
az
├── account          # subscription management
├── acr              # container registries
├── ad               # Microsoft Entra ID (Azure AD)
├── aks              # Azure Kubernetes Service
├── appconfig        # app configuration
├── appservice       # app service plans
├── backup           # azure backup
├── batch            # azure batch
├── bicep            # bicep CLI
├── billing          # billing management
├── cdn              # content delivery networks
├── cognitiveservices # cognitive services
├── config           # CLI configuration
├── configure        # interactive configuration
├── container        # container instances
├── containerapp     # container apps
├── cosmosdb         # cosmos db
├── deployment       # ARM template deployment
├── disk             # managed disks
├── eventgrid        # event grid
├── eventhubs        # event hubs
├── extension        # extension management
├── feature          # resource provider features
├── feedback         # send feedback
├── find             # AI-powered command search
├── functionapp      # function apps
├── graph            # resource graph queries (extension)
├── group            # resource groups
├── identity         # managed identities
├── image            # custom VM images
├── keyvault         # key vault
├── lock             # resource locks
├── login            # authenticate
├── logout           # sign out
├── monitor          # azure monitor
├── mysql            # mysql database
├── network          # networking (vnet, nic, nsg, lb, dns, etc.)
├── policy           # azure policy
├── postgres         # postgresql database
├── provider         # resource providers
├── redis            # redis cache
├── relay            # azure relay
├── rest             # direct REST API calls
├── resource         # generic resource management
├── role             # RBAC role management
├── sf               # service fabric
├── signalr          # signalr service
├── sql              # sql database
├── storage          # storage accounts
├── tag              # resource tags
├── ts               # template specs
├── version          # version info
├── vm               # virtual machines
├── vmss             # virtual machine scale sets
├── webapp           # web apps
└── ...              # 200+ command groups total
```

## Command Group Categories

| Category | Major Groups |
|----------|-------------|
| **Compute** | `vm`, `vmss`, `aks`, `container`, `containerapp`, `functionapp`, `webapp`, `appservice`, `batch`, `image`, `disk`, `sig` |
| **Networking** | `network` (vnet, subnet, nic, nsg, lb, dns, public-ip, application-gateway, etc.), `cdn`, `signalr`, `relay`, `dns-resolver` |
| **Storage & Data** | `storage`, `cosmosdb`, `sql`, `mysql`, `postgres`, `mariadb`, `redis`, `dls`, `netappfiles`, `eventhubs`, `servicebus` |
| **Security & Identity** | `ad`, `role`, `keyvault`, `identity`, `security`, `policy`, `lock`, `attestation` |
| **Management & Operations** | `account`, `group`, `resource`, `provider`, `deployment`, `stack`, `monitor`, `billing`, `policy`, `tag`, `feature` |
| **Containers & DevOps** | `acr`, `aks`, `container`, `containerapp`, `devops` (extension), `pipelines` (extension) |
| **AI & Analytics** | `cognitiveservices`, `search`, `ml` (extension), `kusto` (extension), `synapse`, `databricks` (extension) |
| **Integration** | `eventgrid`, `logic` (extension), `logicapp`, `servicebus`, `eventhubs`, `relay` |

## Common CRUD Operations

Most resource command groups follow a consistent CRUD pattern:

| Operation | Command | Description |
|-----------|---------|-------------|
| Create | `az <resource> create` | Create a new resource (required args: `--name`, `--resource-group`) |
| List | `az <resource> list` | List resources (optional: `--resource-group`) |
| Show | `az <resource> show` | Get details (required: `--name`, `--resource-group`) |
| Delete | `az <resource> delete` | Delete a resource (required: `--name`, `--resource-group`) |
| Update | `az <resource> update` | Update a resource (required: `--name`, `--resource-group`) |
| Wait | `az <resource> wait` | Block until a condition is met |

## Dispatch Model

The az CLI dispatches commands through this sequence:

```
1. Parse command line → identify command path (e.g., "network vnet subnet create")
2. Load command table → only load modules needed for the command
3. Resolve arguments → merge introspection defaults with registered overrides
4. Apply config → resolve defaults from config file and environment variables
5. Authenticate → use cached tokens or refresh via MSAL
6. Execute handler → call the Python handler function
7. Process result → apply --query (JMESPath), then --output formatting
8. Output → write to stdout (results) and stderr (warnings, errors)
```

### Index-Based Loading (Optimized)

Since az CLI 2.x, a command index is used to avoid loading all modules at startup:

1. `CommandIndex.get(args)` looks up the command's top-level name
2. Returns `(modules_list, extensions_list)` needed for that command
3. Only those specific modules are imported — not all 50+ command modules
4. Falls back to full module discovery if the index doesn't contain the command

### Fallback Loading (Full Discovery)

When the index doesn't contain a command:
1. `pkgutil.iter_modules('azure.cli.command_modules')` discovers all installed modules
2. Each module's `COMMAND_LOADER_CLS` is loaded
3. Each loader's `load_command_table()` is called to populate the command table
4. Extensions are discovered and loaded similarly

## Knack Framework Internals

The az CLI is built on Microsoft's **Knack** CLI framework (`pip install knack`), which extends Python's `argparse`.

### Core Knack Components

| Component | Role |
|-----------|------|
| `CLI` | Top-level entry point; holds config, commands loader, and invokes commands |
| `CommandInvoker` | Orchestrates command table loading, parsing, and execution |
| `CLICommandsLoader` | Base loader class — populates `command_table` and registers arguments |
| `CommandGroup` | Context manager that registers commands under a group name |
| `ArgumentsContext` | Context manager that registers argument overrides for a command scope |
| `ArgumentRegistry` | Hierarchical registry — resolves arguments from global to command-specific scope |
| `CLICommand` | Represents a single command with handler, arguments, and help |
| `CLICommandParser` | Extends `argparse` with nested subparser support for command groups |
| `CLIHelp` | Generates help text from parser state and YAML help files |

### Command Table

The command table is an `OrderedDict` mapping command name strings to `CLICommand` instances:

```python
# Example: registering commands in a loader
with CommandGroup(self, 'network vnet', '__main__#{}') as g:
    g.command('create', 'create_vnet')   # → "network vnet create"
    g.command('delete', 'delete_vnet')   # → "network vnet delete"

with CommandGroup(self, 'network vnet subnet', '__main__#{}') as g:
    g.command('create', 'create_subnet') # → "network vnet subnet create"
    g.command('list', 'list_subnets')    # → "network vnet subnet list"
```

The `CommandGroup` context manager concatenates the group name with the command name to form the full command path. It also auto-populates the command group table with all parent groups (e.g., `network`, `network vnet`, `network vnet subnet`).

### Argument Registration

Arguments are registered hierarchically using `ArgumentsContext`:

```python
with ArgumentsContext(self, 'network vnet') as ac:
    ac.argument('vnet_name', options_list=['--name', '-n'])

with ArgumentsContext(self, 'network vnet create') as ac:
    ac.argument('address_prefix', type=str, required=True)
```

The `ArgumentRegistry` resolves arguments by walking from the broadest scope (`""`) to the most specific (full command name), merging `CLIArgumentType` overrides at each level. This allows global defaults to be overridden by command-specific settings.

### Introspection-Based Arguments

Knack auto-extracts arguments from handler function signatures using `inspect.signature()`:

```python
def create_vnet(cmd, resource_group_name, vnet_name, address_prefixes=None):
    """Create a virtual network.
    :param resource_group_name: Name of the resource group
    :param vnet_name: Name of the VNet
    :param address_prefixes: Space-separated IP address prefixes
    """
```

This automatically creates `--resource-group-name`, `--vnet-name`, `--address-prefixes` arguments with help text from the docstring `:param` sections.

### Parser Hierarchy

Nested command groups (e.g., `az network vnet subnet create`) are handled by creating nested subparsers:

```
root parser
  └── subparser: network
       └── subparser: vnet
            └── subparser: subnet
                 └── command: create (with arguments)
                 └── command: delete
                 └── command: list
            └── command: create
            └── command: delete
       └── subparser: nic
       └── subparser: nsg
```

The `CLICommandParser._get_subparser()` method lazily creates subparsers at each level of the command path, caching them in `self.subparsers` dict keyed by tuple paths.

## AzCommandsLoader Extensions

The az CLI extends Knack's base classes:

| Az CLI Class | Knack Base | Azure-Specific Additions |
|--------------|-----------|------------------------|
| `AzCli` | `CLI` | Azure config, cloud profiles, telemetry |
| `MainCommandsLoader` | `CLICommandsLoader` | Module discovery via `pkgutil`, command index, extension loading |
| `AzCliCommand` | `CLICommand` | Local context, confirmation prompts, LRO polling, resource type support |
| `AzCliCommandInvoker` | `CommandInvoker` | Job execution (iterable args), LRO polling, pagination, result transformation |

## Command Module Structure

Each built-in command module lives under `src/azure-cli/azure/cli/command_modules/<name>/`:

| File | Purpose |
|------|---------|
| `__init__.py` | Exports `COMMAND_LOADER_CLS` — the loader class for this module |
| `_command_group.py` | Handler functions (the actual command implementations) |
| `_params.py` | Argument type definitions and validators |
| `_help.py` | YAML help strings for commands |
| `custom.py` | Additional handler functions |

A module must export either `COMMAND_LOADER_CLS` (a class inheriting `AzCommandsLoader`) or a `get_command_loader()` function.

## Events System

Knack provides an event system with hooks at every stage of command execution:

| Event | When |
|-------|------|
| `EVENT_CLI_PRE_EXECUTE` | Before any command runs |
| `EVENT_CLI_SUCCESSFUL_EXECUTE` | After successful execution |
| `EVENT_CLI_POST_EXECUTE` | After execution (success or failure) |
| `EVENT_INVOKER_PRE_CMD_TBL_CREATE` | Before building command table |
| `EVENT_INVOKER_POST_CMD_TBL_CREATE` | After building command table |
| `EVENT_INVOKER_CMD_TBL_LOADED` | After command table is fully loaded |
| `EVENT_INVOKER_PRE_PARSE_ARGS` | Before argument parsing |
| `EVENT_INVOKER_POST_PARSE_ARGS` | After argument parsing |
| `EVENT_INVOKER_TRANSFORM_RESULT` | Before result transformation |
| `EVENT_INVOKER_FILTER_RESULT` | Before result filtering |
| `EVENT_CMDLOADER_LOAD_COMMAND_TABLE` | When loading command table |
| `EVENT_CMDLOADER_LOAD_ARGUMENTS` | When loading arguments |

## Interactive Mode

`az interactive` launches a REPL-style interactive shell:

- Implemented as a core command module that acts as a **launcher** for a separate extension
- If the `interactive` extension is not installed, it is auto-installed
- Uses `prompt_toolkit` for the REPL with tab completion, syntax highlighting, and command history
- Supports `--style` argument with color themes (quiet, purple, default, contrast, etc.)
- Provides inline parameter descriptions and scoped completion

## Top-Level Commands (Non-Group)

These are direct actions, not command groups:

| Command | Description |
|---------|-------------|
| `az login` | Log in to Azure |
| `az logout` | Log out of Azure |
| `az rest` | Invoke a custom REST API request |
| `az find` | AI-powered command recommendations |
| `az feedback` | Send feedback to the Azure CLI team |
| `az upgrade` | Upgrade the CLI and extensions (Preview) |
| `az version` | Show versions of CLI modules and extensions |
| `az survey` | Take the Azure CLI survey |
| `az configure` | Interactive configuration |
| `az init` | Interactive configuration tool (Experimental) |
| `az interactive` | Start interactive mode |

## References

- Knack framework: <https://github.com/microsoft/knack>
- Azure CLI source: <https://github.com/Azure/azure-cli>
- Command reference: <https://learn.microsoft.com/en-us/cli/azure/reference-index>

## Related Skills

- For details on specific service command groups, see [command-groups.md](command-groups.md)
- For global parameters, see [global-parameters.md](global-parameters.md)
- For extension system details, see [extensions.md](extensions.md)
- For authentication, see [auth.md](auth.md)
- For configuration, see [config.md](config.md)
