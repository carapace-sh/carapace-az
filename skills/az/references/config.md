# Configuration System

The Azure CLI configuration system — config file, environment variables, defaults, and cloud profiles.

> **Source of truth**: <https://learn.microsoft.com/en-us/cli/azure/azure-cli-configuration?view=azure-cli-latest>. For global parameters, see [global-parameters.md](global-parameters.md).

## Configuration File

**Location**: `$AZURE_CONFIG_DIR/config` (default: `~/.azure/config` on Linux/macOS, `%USERPROFILE%\.azure\config` on Windows)

**Format**: INI file with `[section]` headers and `key=value` entries:

```ini
[core]
output = table
collect_telemetry = true
only_show_errors = false

[defaults]
group = MyResourceGroup
location = eastus

[storage]
account = mystorageaccount
connection_string = DefaultEndpointsProtocol=https;AccountName=...

[cloud]
name = AzureCloud

[extension]
use_dynamic_install = yes_prompt
```

**INI format rules**:
- Comments start with `#` or `;`
- Section names are case-sensitive
- Key names are not case-sensitive
- Values are strings (booleans are `true`/`false`)
- Whitespace around `=` is trimmed

## Precedence

```
Command-line parameters  (highest)
    ↓
Environment variables (AZURE_{SECTION}_{NAME})
    ↓
Configuration file values  (lowest)
```

## Configuration Commands

| Command | Description |
|---------|-------------|
| `az config set <key>=<value>` | Set a config value |
| `az config get <key>` | Get a config value |
| `az config list` | List all config values |
| `az config unset <key>` | Remove a config value |
| `az configure` | Interactive configuration wizard (GA) |
| `az init` | Interactive configuration tool (Experimental) |

Examples:

```bash
# Set default output format
az config set core.output=table

# Set default resource group
az config set defaults.group=MyResourceGroup

# Set default location
az config set defaults.location=eastus

# Get current output format
az config get core.output

# List all config
az config list --output table
```

## Key Configuration Sections

### `[core]`

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `output` | string | `json` | Default output format |
| `disable_confirm_prompt` | bool | `false` | Turn off confirmation prompts for destructive operations |
| `collect_telemetry` | bool | `true` | Allow anonymous usage data collection |
| `only_show_errors` | bool | `false` | Only show errors, suppress warnings |
| `no_color` | bool | `false` | Disable colored output |
| `enable_broker_on_windows` | bool | `true` (v2.61.0+) | Use WAM for auth on Windows |
| `login_experience_v2` | string | `on` | Subscription selector at login (v2.61.0+) |

### `[defaults]`

Default parameter values for common arguments. When set, you can omit the corresponding `--` flag.

| Key | Corresponds to | Description |
|-----|----------------|-------------|
| `group` | `--resource-group` | Default resource group |
| `location` | `--location` | Default Azure region |
| `web` | `--name` (webapp) | Default web app name |
| `vm` | `--name` (vm) | Default VM name |
| `vmss` | `--name` (vmss) | Default VMSS name |
| `acr` | `--name` (acr) | Default container registry |

```bash
# After setting defaults, these become optional
az config set defaults.group=MyResourceGroup
az config set defaults.location=eastus

# --resource-group and --location can be omitted
az vm create --name MyVm --image Ubuntu2204
```

### `[storage]`

Defaults for `az storage` data-plane commands (blob, file, queue, table operations).

| Key | Description |
|-----|-------------|
| `account` | Default storage account name |
| `key` | Default storage account key |
| `sas_token` | Default SAS token |
| `connection_string` | Default connection string |

### `[cloud]`

| Key | Description |
|-----|-------------|
| `name` | Cloud name: `AzureCloud` (default), `AzureChinaCloud`, `AzureUSGovernment`, `AzureGermanCloud` (deprecated), or custom cloud |

Manage clouds with `az cloud` commands:

```bash
az cloud list --output table
az cloud set --name AzureChinaCloud
az cloud list-available  # List registered clouds
```

### `[extension]`

| Key | Default | Description |
|-----|---------|-------------|
| `use_dynamic_install` | `no` | `no`, `yes_prompt`, `yes_without_prompt` |
| `run_after_dynamic_install` | `true` | Continue running command after dynamic install |
| `index_url` | (Microsoft index) | Custom extension index URL |
| `dir` | `~/.azure/cliextensions` | Extension install directory |

### `[logging]`

| Key | Default | Description |
|-----|---------|-------------|
| `enable_log_file` | `false` | Turn logging to file on/off |
| `log_dir` | `${AZURE_CONFIG_DIR}/logs` | Directory for log files |

## Environment Variables

Each config key can be set as an environment variable using the pattern `AZURE_{SECTION}_{NAME}` (all uppercase, underscores):

| Environment Variable | Corresponds to Config |
|---------------------|----------------------|
| `AZURE_CORE_OUTPUT` | `core.output` |
| `AZURE_CORE_COLLECT_TELEMETRY` | `core.collect_telemetry` |
| `AZURE_CORE_ONLY_SHOW_ERRORS` | `core.only_show_errors` |
| `AZURE_CORE_NO_COLOR` | `core.no_color` |
| `AZURE_DEFAULTS_GROUP` | `defaults.group` |
| `AZURE_DEFAULTS_LOCATION` | `defaults.location` |
| `AZURE_STORAGE_ACCOUNT` | `storage.account` |
| `AZURE_STORAGE_KEY` | `storage.key` |
| `AZURE_STORAGE_SAS_TOKEN` | `storage.sas_token` |
| `AZURE_STORAGE_CONNECTION_STRING` | `storage.connection_string` |
| `AZURE_CLOUD_NAME` | `cloud.name` |
| `AZURE_EXTENSION_USE_DYNAMIC_INSTALL` | `extension.use_dynamic_install` |

### Special Environment Variables

| Variable | Description |
|----------|-------------|
| `AZURE_CONFIG_DIR` | Override the Azure CLI config directory (default: `~/.azure`) |
| `AZURE_EXTENSION_DIR` | Override the extension install directory |
| `AZURE_HTTP_USER_AGENT` | Set a custom User-Agent header |
| `AZURE_ACCESS_TOKEN_FILE` | Provide a pre-obtained access token |

### Using `AZURE_CONFIG_DIR` for Isolation

```bash
# Use separate config directories for different environments
export AZURE_CONFIG_DIR=/tmp/azure-dev
az login --username dev@example.com

export AZURE_CONFIG_DIR=/tmp/azure-prod
az login --username prod@example.com

# Each has its own token cache, config, and active subscription
```

This is also useful for avoiding MSAL token cache conflicts when running multiple `az` scripts concurrently.

## Cloud Profiles

The Azure CLI supports multiple cloud environments (sovereign clouds):

| Cloud Name | Description |
|------------|-------------|
| `AzureCloud` | Public Azure (default) |
| `AzureChinaCloud` | Azure China (21Vianet) |
| `AzureUSGovernment` | Azure US Government |
| `AzureGermanCloud` | Azure Germany (deprecated) |

Each cloud has its own set of API endpoints and management URLs. The `az cloud` commands manage cloud registrations:

```bash
# List available clouds
az cloud list --output table

# Switch to a different cloud
az cloud set --name AzureUSGovernment

# Register a custom cloud
az cloud register --name MyCloud --endpoint-resource-manager "https://management.mycompany.com"
```

## `az configure` vs `az config set` vs `az init`

| Command | Status | Description |
|---------|--------|-------------|
| `az configure` | GA | Interactive wizard — walks through common settings (output format, telemetry, file logging) |
| `az config set` | GA | Direct key=value setting for any config property |
| `az init` | Experimental | AI-assisted interactive configuration tool with scenario-based setup |

## References

- Configuration docs: <https://learn.microsoft.com/en-us/cli/azure/azure-cli-configuration>
- Cloud management: <https://learn.microsoft.com/en-us/cli/azure/cloud>
- Environment variables: <https://learn.microsoft.com/en-us/cli/azure/azure-cli-configuration#configure-environment-variables>

## Related Skills

- For global parameters that interact with config, see [global-parameters.md](global-parameters.md)
- For authentication and token cache, see [auth.md](auth.md)
- For extension configuration, see [extensions.md](extensions.md)
- For `AZURE_CONFIG_DIR` gotchas, see [gotchas.md](gotchas.md)
