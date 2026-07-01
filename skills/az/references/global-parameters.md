# Global Parameters

The parameters available on (nearly) every Azure CLI command.

> **Source of truth**: <https://learn.microsoft.com/en-us/cli/azure/azure-cli-global-parameters?view=azure-cli-latest>. For output formatting and queries, see [output-formatting.md](output-formatting.md). For configuration, see [config.md](config.md).

## Parameter Precedence

```
Command-line parameters  (highest)
    ↓
Environment variables (AZURE_{SECTION}_{NAME})
    ↓
Configuration file values  (lowest)
```

## Primary Global Parameters

| Parameter | Short | Description | Default |
|-----------|-------|-------------|---------|
| `--help` | `-h` | Show help message and exit | |
| `--output` | `-o` | Output format | `json` |
| `--query` | | JMESPath query string | |
| `--subscription` | | Name or ID of subscription | active subscription |
| `--verbose` | | Increase logging verbosity | `False` |
| `--debug` | | Show all debug logs | `False` |
| `--only-show-errors` | | Only show errors, suppress warnings | `False` |
| `--version` | `-v` | Show version (top-level only) | |

## Output Format Values

| Value | Description |
|-------|-------------|
| `json` | JSON string (default) |
| `jsonc` | Colorized JSON |
| `table` | ASCII table with keys as column headings; nested objects not included |
| `tsv` | Tab-separated values, no keys, no nested objects |
| `yaml` | YAML format |
| `yamlc` | Colorized YAML |
| `none` | No output other than errors and warnings |

Change default globally: `az config set core.output=table`

For JMESPath query details, see [output-formatting.md](output-formatting.md).

## `--subscription` Behavior

The `--subscription` parameter overrides the active subscription for a single command without changing the persistent default:

```bash
# Use a specific subscription for this command only
az vm list --subscription "My Dev Subscription" --output table

# The active subscription is unchanged for subsequent commands
az vm list --output table  # still uses the default subscription
```

Accepts both subscription ID (GUID) and subscription name. If the subscription belongs to a different tenant, the active tenant context also changes for that command.

## `--debug` vs `--verbose`

| Level | What It Shows |
|-------|---------------|
| (default) | Results and warnings |
| `--verbose` | Results, warnings, and informational messages (e.g., request URLs, status) |
| `--debug` | Everything from `--verbose` plus full debug logs (headers, request/response bodies, MSAL details, retry attempts) |

`--debug` is useful for troubleshooting API calls and authentication issues. It outputs the full HTTP request and response including headers (with authorization tokens redacted).

## `--only-show-errors`

Suppresses warnings, preview notices, and deprecation notices — only errors are shown. Useful in scripts where you want clean output:

```bash
az vm list --only-show-errors --output table
```

Can be set globally: `az config set core.only_show_errors=true`

## `--query` Interaction with `--output`

The `--query` parameter applies JMESPath queries **before** output formatting. The result type determines how output formats behave:

| Query Result Type | `json` | `table` | `tsv` | `yaml` |
|-----------------|--------|---------|-------|--------|
| Array of objects | Array | Table (keys as columns) | One row per object | YAML sequence |
| Single object | Object | Table (properties as rows) | Single row | YAML mapping |
| Array of primitives | Array | Table (Column1) | One per line | YAML sequence |
| Single primitive | Value | Table (single cell) | Single value | Scalar |

For table output, certain keys (`id`, `type`, `etag`) are automatically filtered out. To display them, rename via multiselect hash:

```bash
az vm show -g MyGroup -n MyVm --query "{resourceId:id}" --output table
```

## Environment Variable Equivalents

Some global parameters have environment variable equivalents:

| Parameter | Environment Variable |
|-----------|---------------------|
| `--output` | `AZURE_CORE_OUTPUT` |
| `--subscription` | `AZURE_SUBSCRIPTION` (not officially documented) |
| `--only-show-errors` | `AZURE_CORE_ONLY_SHOW_ERRORS` |
| `--debug` | `AZURE_DEBUG` (not official) |

For the full environment variable list, see [config.md](config.md).

## Shell-Specific Considerations

### Bash

```bash
# Double quotes for --query, single quotes for string literals inside
az vm list --query "[?location=='eastus'].name" --output table

# Backticks in JMESPath need escaping
az vm list --query "[?diskSizeGb >=\`50\`].name"
```

### PowerShell

```powershell
# Backticks need doubling (PowerShell escape char conflicts with JMESPath)
az vm list --query "[?diskSizeGb >=``50``].name"

# Single-quoted strings work, but double single quotes for literals
az vm list --query '[?location==''eastus''].name'
```

### Cmd

```cmd
REM Backticks are literal - no escaping needed
az vm list --query "[?diskSizeGb >=`50`].name"

REM Use 'call' prefix in batch scripts
call az vm list --query "[?location=='eastus'].name"
```

For detailed quoting rules, see [output-formatting.md](output-formatting.md) and [gotchas.md](gotchas.md).

## References

- Global parameters: <https://learn.microsoft.com/en-us/cli/azure/azure-cli-global-parameters>
- Output formatting: <https://learn.microsoft.com/en-us/cli/azure/format-output-azure-cli>
- JMESPath queries: <https://learn.microsoft.com/en-us/cli/azure/query-azure-cli>

## Related Skills

- For JMESPath query syntax, see [output-formatting.md](output-formatting.md)
- For configuration and environment variables, see [config.md](config.md)
- For subscription management, see [auth.md](auth.md)
- For shell quoting gotchas, see [gotchas.md](gotchas.md)
