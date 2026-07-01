# Output Formatting and JMESPath Queries

How `--output` and `--query` work in the Azure CLI to format and transform command results.

> **Source of truth**: <https://learn.microsoft.com/en-us/cli/azure/format-output-azure-cli> and <https://learn.microsoft.com/en-us/cli/azure/query-azure-cli>. For global parameter details, see [global-parameters.md](global-parameters.md).

## Processing Pipeline

```
Command handler result (Python object)
  → JSON serialization
    → --query (JMESPath query applied to JSON)
      → --output formatting
        → stdout
```

The `--query` is always applied to the JSON representation **before** output formatting. The output format only affects how the query result is displayed.

## Output Formats

| Format | Description | Use Case |
|--------|-------------|----------|
| `json` | Pretty-printed JSON string (default) | Full data, piping to `jq` |
| `jsonc` | Colorized JSON | Interactive inspection |
| `table` | ASCII table with keys as column headings | Quick visual scanning |
| `tsv` | Tab-separated values, no keys, no nesting | Scripting, variable assignment |
| `yaml` | YAML format | Human-readable, config-like |
| `yamlc` | Colorized YAML | Interactive inspection |
| `none` | No output (only errors and warnings) | Suppress output in scripts |

Change default: `az config set core.output=table`

### Table Output Behavior

- Nested objects and arrays are **not shown** in table output
- Certain keys are automatically filtered out: `id`, `type`, `etag`
- To show filtered keys, rename them in a multiselect hash:

```bash
# 'id' is filtered out of table output
az vm show -g MyGroup -n MyVm --query "{resourceId:id}" --output table
```

- If you use a multiselect **list** (square brackets), columns are named `Column1`, `Column2`, etc.

```bash
az vm list --query "[*].[name, hardwareProfile.vmSize]" --output table
# Column1 | Column2
```

### TSV Output

- No keys, no quotes, no JSON structure — just raw values
- Tab-separated columns, newline-separated rows
- **Column ordering is alphabetic by key** — use `--query` to force ordering:

```bash
# WARNING: column order is alphabetic, not insertion order
az vm list --output tsv  # columns may not be in expected order

# Force ordering with multiselect list
az vm list --query "[].[name, location, resourceGroup]" --output tsv
```

- For single values, removes JSON quotes:

```bash
# WITHOUT tsv - includes JSON quotes
USER=$(az vm show -g MyGroup -n MyVm --query "osProfile.adminUsername")
echo $USER  # "azureuser"  ← includes quotes!

# WITH tsv - clean value
USER=$(az vm show -g MyGroup -n MyVm --query "osProfile.adminUsername" --output tsv)
echo $USER  # azureuser  ← clean
```

## JMESPath Query Syntax

JMESPath is a query language for JSON. Reference: <https://jmespath.org/specification.html>

### Property Access (`.`)

```bash
# Nested property access (case-sensitive)
az vm show -g MyGroup -n MyVm --query "osProfile.linuxConfiguration.ssh.publicKeys"
```

### Multiselect List (`[expr1, expr2, ...]`)

Returns an array of values (no keys):

```bash
az vm show -g MyGroup -n MyVm \
  --query "[name, osProfile.adminUsername, osProfile.linuxConfiguration.ssh.publicKeys[0].keyData]"
# Output: ["MyVm", "azureuser", "ssh-rsa AAAAB3..."]
```

### Multiselect Hash (`{key:expr, ...}`)

Returns a dictionary with custom key names:

```bash
az vm show -g MyGroup -n MyVm \
  --query "{VMName:name, admin:osProfile.adminUsername, sshKey:osProfile.linuxConfiguration.ssh.publicKeys[0].keyData}"
# Output: {"VMName": "MyVm", "admin": "azureuser", "sshKey": "ssh-rsa AAAAB3..."}
```

### Flattening (`[]`)

Applies subsequent expressions to **each element** in an array:

```bash
# Flatten top-level result array
az vm list -g MyGroup --query "[].{Name:name, OS:storageProfile.osDisk.osType}"

# Flatten a nested array
az vm show -g MyGroup -n MyVm --query "osProfile.linuxConfiguration.ssh.publicKeys[].keyData"
```

### Filtering (`[?predicate]`)

Filters array elements where the predicate is true:

```bash
# Boolean property (true)
az account list --query "[?isDefault].name"

# Boolean property (false - negation)
az account list --query "[?!isDefault].name"

# String comparison
az vm list -g MyGroup --query "[?storageProfile.osDisk.osType=='Linux'].{Name:name, admin:osProfile.adminUsername}"

# Numeric comparison (backtick-escape numbers)
az vm list -g MyGroup --query "[?storageProfile.osDisk.diskSizeGb >=\`50\`].{Name:name, DiskSize:storageProfile.osDisk.diskSizeGb}"
```

### Logical Operators

```bash
# AND
az vm list --query "[?location=='eastus' && hardwareProfile.vmSize=='Standard_DS2_v2']"

# OR
az vm list --query "[?tags.env == 'test' || tags.env == 'dev']"
```

### Pipe Expressions (`|`)

Pass intermediate results to the next expression:

```bash
# Project properties, then filter the projected result
az vm list -g MyGroup \
  --query "[].{Name:name, Storage:storageProfile.osDisk.managedDisk.storageAccountType} | [? contains(Storage,'SSD')]"
```

### Slicing Arrays

```bash
# First element
--query "[0]"

# Last element
--query "[-1]"

# First two elements
--query "[:2]"

# Elements from index 2 onwards
--query "[2:]"

# Every other element
--query "[::2]"

# Reverse the array
--query "[::-1]"
```

## JMESPath Functions

| Function | Purpose | Example |
|----------|---------|---------|
| `contains(string, substring)` | Filter by substring | `[?contains(storageProfile.osDisk.managedDisk.storageAccountType,'SSD')]` |
| `starts_with(string, prefix)` | Filter by prefix | `[?starts_with(name, 'vm-')]` |
| `ends_with(string, suffix)` | Filter by suffix | `[?ends_with(name, '-prod')]` |
| `sort_by(array, &expr)` | Sort objects by property | `sort_by([], &Size)` |
| `sort(@)` | Sort simple values | `[].name \| sort(@)` |
| `reverse(array)` | Reverse order | `reverse(sort_by([], &Size))` |
| `length(array)` | Count elements | `length([?osType=='Linux'])` |
| `min_by(array, &expr)` | Min by property | `min_by([], &diskSizeGb)` |
| `max_by(array, &expr)` | Max by property | `max_by([], &diskSizeGb)` |
| `to_string(value)` | Convert to string | `sort_by(..., &to_string(dataAction))` (boolean sort workaround) |

The `&` prefix on sort expressions tells JMESPath to **defer evaluation** — evaluated per element during sorting.

## Common Query Patterns

### Filtering by Property

```bash
# String equality
--query "[?state=='Running']"

# String contains
--query "[?contains(name, 'prod')]"

# Numeric comparison
--query "[?diskSizeGb >=\`50\`]"

# Boolean true
--query "[?isDefault]"

# Boolean false
--query "[?!isDefault]"
# or
--query "[?isDefault == \`false\`]"
```

### Projecting Specific Properties

```bash
# Single property from each array element
--query "[].id"

# Multiple properties with custom keys
--query "[].{VMName:name, Location:location, Size:hardwareProfile.vmSize}"
```

### Sorting

```bash
# Sort simple values
--query "[].name | sort(@)"

# Sort objects by a property
--query "sort_by([].{Name:name, Size:storageProfile.osDisk.diskSizeGb}, &Size)"

# Sort descending
--query "reverse(sort_by([], &storageProfile.osDisk.diskSizeGb))"
```

### Nested Array Projection

```bash
# Project from nested arrays
az keyvault show -n MyKv -g MyGroup \
  --query "properties.accessPolicies[*].{objectId:objectId, permissions:permissions}[]"

# Filter nested arrays
az keyvault show -n MyKv -g MyGroup \
  --query "properties.accessPolicies[*].{objectId:objectId, permissions:permissions}[? contains(permissions.secrets, 'Delete')]"
```

### Using Shell Variables in Queries

```bash
# Bash
IP="20.127"
az network public-ip list --query "[?ipAddress!=null]|[?contains(ipAddress, '$IP')].[id]" --output tsv

# PowerShell
$IP = "20.127"
az network public-ip list --query "[?ipAddress!=null]|[?contains(ipAddress, '$IP')].[id]" --output tsv

# Cmd
set IP=20.127
az network public-ip list --query "[?ipAddress!=null]|[?contains(ipAddress, '%IP%')].[id]" --output tsv
```

## Shell Quoting Reference

The #1 source of errors in Azure CLI queries is shell quoting. Different shells handle quoting very differently.

### Bash

- Use **double quotes** around the query: `--query "expression"`
- Single quotes inside for string literals: `--query "[?state=='Running']"`
- Backticks need escaping with backslash: `` \` ``

```bash
# Good
az vm list --query "[?state=='Running'].name"

# Numeric - escape backticks
az vm list --query "[?diskSizeGb >=\`50\`].name"
```

### PowerShell

PowerShell's escape character (`` ` ``) collides with JMESPath's literal escape:

- **Double backticks** for escaped values in PowerShell 7+:

```powershell
az vm list --query "[?diskSizeGb >=``50``].{Name:name}"
az account list --query "[?isDefault ==``false``].name"
```

- Single-quoted strings with doubled single quotes:

```powershell
az vm list --query '[?state==''Running''].name'
```

### Cmd (Command Prompt)

- Use **double quotes** for the query
- Backticks are **literal** — no escaping needed
- Prefix with `call` in batch scripts

```cmd
az vm list --query "[?diskSizeGb >=`50`].name"
call az account list --query "[?isDefault].name"
```

### Quoting Summary

| Construct | Bash | PowerShell | Cmd |
|-----------|------|------------|-----|
| Query wrapper | `"..."` | `"..."` or `'...'` | `"..."` |
| String `'Linux'` | `'Linux'` | `'Linux'` or `''Linux''` | `'Linux'` |
| Number `` `50` `` | `` \`50\` `` | ` ``50`` ` | `` `50` `` |
| Boolean `` `false` `` | `` \`false\` `` | ` ``false`` ` | `` `false` `` |
| Shell variable | `'$VAR'` | `'$VAR'` | `%VAR%` |
| Batch prefix | N/A | N/A | `call` before `az` |

## Query Result Type → Output Format Interaction

| Query Result Type | `json` | `table` | `tsv` | `yaml` |
|-----------------|--------|---------|-------|--------|
| Array of objects | JSON array | Table (keys as columns) | One row per object | YAML sequence |
| Single object | JSON object | Table (properties as rows) | Single row | YAML mapping |
| Array of primitives | JSON array | Table (Column1) | One per line | YAML sequence |
| Single primitive | JSON value | Table (single cell) | Single value | YAML scalar |

## Edge Cases

### Incorrect Quotes in Filter Predicates

JMESPath strings use **single quotes** (`'`) or escape characters (`` ` ``). Using double quotes inside a filter predicate produces **empty output**:

```bash
# WRONG - double quotes inside predicate → empty output
az vm list --query "[?osType==\"Linux\"]"  # returns []

# CORRECT - single quotes
az vm list --query "[?osType=='Linux']"
```

### Spaces in Column Names

If using spaces in multiselect hash keys (e.g., `{VM Name:name}`), quoting rules change in both Bash and PowerShell. Avoid spaces for simplicity — use `VMName` or `vm_name` instead.

### Boolean Sorting Workaround

JMESPath `sort_by` doesn't natively support boolean values. Use `to_string()` as a workaround:

```bash
--query "sort_by(resourceTypes[].operations[].{action:name, dataAction:isDataAction}, &action) | sort_by(@, &to_string(dataAction))"
```

## References

- Output formatting: <https://learn.microsoft.com/en-us/cli/azure/format-output-azure-cli>
- JMESPath queries: <https://learn.microsoft.com/en-us/cli/azure/query-azure-cli>
- JMESPath specification: <https://jmespath.org/specification.html>
- JMESPath tutorial: <https://jmespath.org/tutorial.html>

## Related Skills

- For global parameters, see [global-parameters.md](global-parameters.md)
- For shell quoting gotchas, see [gotchas.md](gotchas.md)
- For configuration of default output format, see [config.md](config.md)
