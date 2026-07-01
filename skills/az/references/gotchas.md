# Gotchas and Pitfalls

Common pitfalls, edge cases, and non-obvious behaviors in the Azure CLI.

> **Source of truth**: <https://learn.microsoft.com/en-us/cli/azure/use-azure-cli-successfully-tips> and <https://learn.microsoft.com/en-us/cli/azure/use-azure-cli-successfully-troubleshooting>. For detailed parameter behavior, see [global-parameters.md](global-parameters.md).

## Shell Quoting Differences

The #1 source of errors. Azure CLI docs are primarily written and tested in **Bash**. Commands with `--query`, JSON strings, or special characters often break when copied to PowerShell or cmd.exe.

```bash
# Bash - works
az vm list --query "[?location=='eastus'].name"

# PowerShell - backtick collision (need double backticks)
az vm list --query "[?diskSizeGb >=``50``].name"

# Cmd - backticks are literal
az vm list --query "[?diskSizeGb >=`50`].name"
```

See [output-formatting.md](output-formatting.md#shell-quoting-reference) for the full quoting reference table.

**Key rules**:
- JMESPath strings use **single quotes** (`'`), never double quotes inside predicates
- PowerShell's `` ` `` escape character collides with JMESPath `` ` `` literal escape
- Cmd batch scripts need `call` prefix: `call az ...`

## Default Output is JSON

If you assign command output to variables without `--output tsv`, the value includes JSON formatting:

```bash
# BAD - includes JSON quotes and brackets
SUB_ID=$(az account show --query id)
echo $SUB_ID  # "abc-123"  ← includes quotes!

# GOOD - use tsv for clean values
SUB_ID=$(az account show --query id --output tsv)
echo $SUB_ID  # abc-123  ← clean
```

## TSV Column Ordering is NOT Guaranteed

The order of columns in TSV output is **alphabetic by key**, which can change unexpectedly:

```bash
# WARNING: columns are alphabetic, not in the order you might expect
az vm list --output tsv
# Columns might be: location, name, resourceGroup, type (alphabetical)

# FIX: use --query to force column ordering
az vm list --query "[].[name, location, resourceGroup, type]" --output tsv
```

## MSAL Token Cache Concurrency

If running multiple `az` scripts simultaneously, they can write to the same MSAL token cache file (`~/.azure/msal_token_cache.json`) and conflict, causing authentication failures.

**Fix**: Use `AZURE_CONFIG_DIR` to separate directories for each concurrent script:

```bash
# Script 1
export AZURE_CONFIG_DIR=/tmp/azure-session-1
az login --service-principal -u app1 -p secret1 --tenant tenant1

# Script 2 (in a different process)
export AZURE_CONFIG_DIR=/tmp/azure-session-2
az login --service-principal -u app2 -p secret2 --tenant tenant2
```

Each directory gets its own token cache, config file, and active subscription context.

## Extensions Don't Auto-Update

Unlike core CLI commands, extensions must be updated **manually**:

```bash
# Extensions are not updated automatically alongside the CLI core
# You must update each extension separately
az extension update --name ml
az extension update --name devops

# Or update all at once (if supported)
az extension list --query "[].name" --output tsv | xargs -I {} az extension update --name {}
```

`az upgrade` updates the CLI core and can offer to update extensions, but it doesn't always do so automatically.

## `--no-wait` is Silent

Commands with `--no-wait` return immediately with **no output**. To block until the operation completes, use `az <resource> wait`:

```bash
# Create VM without waiting
az vm create --name MyVm --resource-group MyGroup --image Ubuntu2204 --no-wait

# Block until the VM is created
az vm wait --created --name MyVm --resource-group MyGroup

# Or wait for an existing resource by ID
az vm wait --updated --ids /subscriptions/.../providers/Microsoft.Compute/virtualMachines/MyVm
```

## `az rest` — Universal REST API Fallback

When no CLI command exists for an Azure resource, use `az rest` to call the Azure REST API directly. It is automatically authenticated using the active subscription's credentials:

```bash
# GET request
az rest --method get \
  --url https://management.azure.com/subscriptions/{id}/resourcegroups?api-version=2021-04-01

# POST request with body
az rest --method post \
  --url https://management.azure.com/.../restart?api-version=2021-04-01 \
  --body @body.json

# Use --uri to prepend the ARM base endpoint (subscription ID still needs to be filled in)
az rest --method get --uri /subscriptions/{id}/providers/Microsoft.Resources/resources?api-version=2021-04-01
```

**Key flags**:
- `--method`: `get`, `post`, `put`, `patch`, `delete`
- `--url` / `--uri`: Full URL or relative path
- `--body`: Request body (string or `@file.json`)
- `--resource`: Override the resource endpoint (default: `https://management.azure.com/`)
- `--scope`: Override the scope for token acquisition

## `az vm list` vs `az vm list --show-details`

By default, `az vm list` does **not** show power states. You must pass `--show-details` to see them, which triggers an additional API call **per VM** — this can be slow for large numbers of VMs:

```bash
# Fast - no power state
az vm list --output table

# Slow - includes power state (extra API call per VM)
az vm list --show-details --output table
```

## Configuration File is Strict INI

The Azure CLI config file uses INI format with specific rules:

- Comments start with `#` or `;`
- **Section names are case-sensitive**: `[core]` ≠ `[Core]`
- Key names are **not** case-sensitive: `output = table` = `Output = table`
- Values are strings: booleans are `true`/`false` (lowercase)
- No quoting needed for string values
- Whitespace around `=` is trimmed

```ini
# This works
[core]
output = table

# This doesn't (wrong section name case)
[Core]
output = table  # will be ignored
```

## Telemetry is ON by Default

The Azure CLI collects anonymous usage data by default. Disable it:

```bash
az config set core.collect_telemetry=false
```

Or via environment variable: `AZURE_CORE_COLLECT_TELEMETRY=false`

## Core + Extension Dual Status

Many command groups (like `az vm`, `az storage`, `az network`, `az sql`, `az aks`, etc.) exist in **both** core and as extensions. When an extension version is installed that provides the same commands:

- The extension's commands **override** the core commands (with a warning)
- This is intentional — extensions can ship newer/fixed versions of commands
- Removing the extension restores the core version

```
The installed extension 'storage-preview' adds or overwrites some 'az storage' commands. It can be removed using 'az extension remove --name storage-preview'.
```

## `az upgrade` is Still Preview

Despite being available for years, `az upgrade` remains in **Preview** status. It may not work in all environments, especially when the CLI was installed via a package manager (apt, yum, brew, etc.):

```bash
# May not work if installed via package manager
az upgrade

# For package-manager installations, use the package manager to update
# apt: sudo apt update && sudo apt upgrade azure-cli
# yum: sudo yum update azure-cli
# brew: brew upgrade azure-cli
```

## `az spring` is Deprecated

The `az spring` command group (Azure Spring Apps) is marked as **deprecated**. Use the replacement service or extension when available.

## Many "Core" Commands are Actually Extensions

Commands that users often assume are core but are actually extensions:

| Command | Extension |
|---------|-----------|
| `az ml` | `ml` |
| `az devops` | `devops` |
| `az graph` | `graph` |
| `az boards` | `devops` |
| `az repos` | `devops` |
| `az pipelines` | `devops` |
| `az sentinel` | `sentinel` |
| `az ssh` | `ssh` |
| `az terraform` | `terraform` |
| `az databricks` | `databricks` |
| `az datafactory` | `datafactory` |
| `az k8s-extension` | `k8s-extension` |
| `az connectedk8s` | `connectedk8s` |

Without the extension installed, the command won't be found (unless dynamic install is enabled).

## Subscription Cache Staleness

If permissions to a new subscription are granted while your terminal is open, you may see a "subscription doesn't exist" error. Fix by:

1. Closing and reopening the terminal, **or**
2. `az logout` then `az login`, **or**
3. `az account clear` then `az login`

## `az group delete` is Destructive and Async

`az group delete` deletes **all** resources in a resource group. It's asynchronous by default:

```bash
# Deletes everything in MyResourceGroup - returns immediately
az group delete --name MyResourceGroup --no-wait

# To wait for completion
az group delete --name MyResourceGroup  # waits by default (no --no-wait)

# Check deletion status
az group exists --name MyResourceGroup  # returns false when fully deleted
```

## `az storage` Auth Defaults

The `az storage` data-plane commands (blob, file, queue, table) use a separate auth system from the management-plane commands. They look for defaults in the `[storage]` config section:

```ini
[storage]
account = mystorageaccount
key = base64encodedkey==
```

If not set, you must provide `--account-name` and either `--account-key`, `--sas-token`, or `--connection-string` on each data-plane command. Alternatively, use `--auth-mode login` for Microsoft Entra ID-based access:

```bash
az storage blob list --container-name mycontainer --account-name mystorageaccount --auth-mode login
```

## WAM Issues on Windows

If WAM (Web Account Manager) causes authentication issues on Windows:

```bash
# Disable WAM
az config set core.enable_broker_on_windows=false

# Clear existing auth state
az account clear

# Login again (will use browser-based flow)
az login
```

Common WAM error: "User cancelled the Accounts Control Operation" — the user dismissed the WAM dialog or it timed out.

## References

- Tips for success: <https://learn.microsoft.com/en-us/cli/azure/use-azure-cli-successfully-tips>
- Troubleshooting: <https://learn.microsoft.com/en-us/cli/azure/use-azure-cli-successfully-troubleshooting>
- Quoting issues: <https://github.com/Azure/azure-cli/blob/dev/doc/quoting-issues.md>

## Related Skills

- For shell quoting reference, see [output-formatting.md](output-formatting.md#shell-quoting-reference)
- For global parameters, see [global-parameters.md](global-parameters.md)
- For configuration file format, see [config.md](config.md)
- For extension management, see [extensions.md](extensions.md)
- For authentication and WAM details, see [auth.md](auth.md)
