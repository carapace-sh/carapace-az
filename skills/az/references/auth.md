# Authentication and Identity

The Azure CLI authentication system, subscription management, Microsoft Entra ID, and RBAC.

> **Source of truth**: <https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli?view=azure-cli-latest>. For account/subscription commands, see <https://learn.microsoft.com/en-us/cli/azure/account>.

## Authentication Methods

| Method | Command | Best For |
|--------|---------|----------|
| Interactive (browser) | `az login` | Local development, learning |
| Service principal | `az login --service-principal` | Automation, CI/CD, scripts |
| Managed identity | `az login --identity` | Azure-hosted apps (VMs, Functions, etc.) |
| Azure Cloud Shell | (auto-logged in) | Easiest start |

## Interactive Login

```bash
az login
```

- Opens a browser for OAuth2 authorization code flow
- Falls back to **device code flow** if no browser is available
- Force device code: `az login --use-device-code`
- Specify tenant: `az login --tenant <tenant-id>`
- With username/password (not recommended for MFA): `az login --user <username> --password <password>`
- With custom scope: `az login --scope https://management.core.windows.net/.default`

### Subscription Selector (v2.61.0+)

After `az login`, if you have access to multiple subscriptions, you're prompted to select one interactively. The previously selected default is marked with `*`. Pressing Enter keeps the default.

Disable: `az config set core.login_experience_v2=off`

## WAM (Web Account Manager) on Windows

Since Azure CLI v2.61.0 (May 2024), **WAM is the default authentication method on Windows**.

| Aspect | Details |
|--------|---------|
| Platform | Windows 10+, Windows Server 2019+ |
| Benefits | Enhanced security, Windows Hello, FIDO key support, Conditional Access, SSO |
| Disable | `az config set core.enable_broker_on_windows=false` then `az account clear` and `az login` |
| Error | "User cancelled the Accounts Control Operation" — user dismissed the WAM dialog |

WAM is a Windows authentication broker that manages authentication handshakes and token maintenance, providing better security than browser-based flow.

## Service Principal Authentication

Best for automation, scripts, and CI/CD pipelines.

### With Client Secret

```bash
az login --service-principal \
  --username <app-id> \
  --password <client-secret> \
  --tenant <tenant-id>
```

### With Certificate

```bash
az login --service-principal \
  --username <app-id> \
  --certificate /path/to/cert.pem \
  --tenant <tenant-id>
```

The PEM file must contain the certificate appended to the private key.

### Creating a Service Principal

```bash
# Create with RBAC role
az ad sp create-for-rbac --name MyServicePrincipal --role Contributor --scopes /subscriptions/{id}/resourceGroups/{rg}

# Output includes appId, password, tenant
```

## Managed Identity Authentication

No credentials to manage — available on Azure-hosted resources.

### System-Assigned

```bash
az login --identity
```

### User-Assigned

```bash
az login --identity --client-id <client-id>
az login --identity --object-id <object-id>
az login --identity --resource-id <resource-id>
```

## MSAL Token Cache

The Azure CLI uses **MSAL** (Microsoft Authentication Library) for token management:

| Aspect | Details |
|--------|---------|
| Token cache location | `~/.azure/msal_token_cache.json` |
| Token refresh | Automatic — refresh tokens exchanged as needed |
| Token validity | Access tokens valid 5-60 minutes |
| Concurrency issue | Multiple concurrent `az` processes can conflict writing to the same token cache file |
| Fix for concurrency | Set `AZURE_CONFIG_DIR` to separate directories for each concurrent script |

## `az account` — Subscription Management

### `az account list`

List all subscriptions for the logged-in account.

```bash
# Default: only 'Enabled' subscriptions from current cloud
az account list --output table

# All subscriptions, including disabled
az account list --all --output table

# Refresh from server
az account list --refresh --output table

# Get default subscription
az account list --query "[?isDefault]" --output table

# Search by name
az account list --query "[?contains(name,'dev')].{Name:name, ID:id, Tenant:tenantId}" --output table
```

### `az account set`

Set the active subscription for all subsequent commands:

```bash
az account set --subscription "<subscription-id-or-name>"
```

If the subscription is in a different tenant, the active tenant also changes. For per-command `--subscription` override, see [global-parameters.md](global-parameters.md).

### `az account show`

Show details of the active subscription (or a specified one):

```bash
az account show --output table
az account show --subscription "<subscription-id>"
```

### `az account get-access-token`

Get a token for accessing Azure services:

```bash
# For active subscription
az account get-access-token

# For specific subscription
az account get-access-token --subscription "<subscription-id>"

# For specific tenant
az account get-access-token --tenant <tenant-id>

# For MS Graph API
az account get-access-token --resource-type ms-graph
```

Returns JSON with: `accessToken`, `expiresOn` (local datetime), `expires_on` (POSIX/UTC timestamp), `subscription`, `tenant`, `tokenType`.

**Use `expires_on` (POSIX/UTC) for downstream applications** — avoids Daylight Saving Time "fold" issues with `expiresOn`.

Resource types: `aad-graph`, `arm`, `batch`, `data-lake`, `media`, `ms-graph`, `oss-rdbms`.

### `az account clear`

Clears all subscriptions from the CLI's local cache. Not the same as logout — but after clearing, you must `az login` again before running any command.

### `az account list-locations`

```bash
az account list-locations --output table
az account list-locations --query "[?metadata.regionType=='Physical'].name" --output tsv
```

### `az account management-group`

Manage management groups for organizing subscriptions:

```bash
az account management-group create --name MyMG --display-name "My Management Group"
az account management-group list --output table
az account management-group subscription add --name MyMG --subscription <sub-id>
az account management-group subscription remove --name MyMG --subscription <sub-id>
```

## Multi-Tenant Authentication

A tenant is an instance of Microsoft Entra ID for a single organization.

```bash
# Sign in to a specific tenant
az login --tenant <tenant-id>

# Sign in with tenant domain
az login --tenant mycompany.onmicrosoft.com

# Get token for a specific tenant
az account get-access-token --tenant <tenant-id>

# List tenants you have access to (experimental)
az account tenant list
```

Changing the active subscription to one in a different tenant automatically changes the tenant context. If `az login` fails with "Authentication failed against tenant", use `--tenant` to specify the target explicitly.

## MFA Impact (October 2025)

Starting October 2025, Microsoft requires **MFA for Azure CLI** for user identities.

| Identity Type | MFA Required? |
|--------------|---------------|
| User (interactive login) | Yes |
| User (username/password) | Yes (will fail without MFA — migrate to workload identities) |
| Service principal | No |
| Managed identity | No |

Use `--claims-challenge` parameter for MFA scenarios:

```bash
az login --tenant <tenant-id> --claims-challenge "<token>"
```

## `az ad` — Microsoft Entra ID

Manage Microsoft Entra ID (formerly Azure AD) entities via Microsoft Graph API.

### Users

| Command | Description |
|---------|-------------|
| `az ad user create` | Create a user |
| `az ad user delete` | Delete a user |
| `az ad user list` | List users |
| `az ad user show` | Get user details |
| `az ad user update` | Update a user |
| `az ad user get-member-groups` | Get groups for a user |

### Groups

| Command | Description |
|---------|-------------|
| `az ad group create` | Create a group |
| `az ad group delete` | Delete a group |
| `az ad group list` | List groups |
| `az ad group show` | Get group details |
| `az ad group member add/list/remove/check` | Manage group members |
| `az ad group owner add/list/remove` | Manage group owners |
| `az ad group get-member-groups` | Get parent groups |

### Service Principals

| Command | Description |
|---------|-------------|
| `az ad sp create --id` | Create SP from existing application |
| `az ad sp create-for-rbac` | Create app + SP + RBAC role in one step |
| `az ad sp delete --id` | Delete an SP |
| `az ad sp list` | List SPs (default: first 100; `--all` for all) |
| `az ad sp show --id` | Get SP details (by appId, objectId, or identifier URI) |
| `az ad sp update --id` | Update an SP |
| `az ad sp credential reset` | Reset password or certificate credentials |
| `az ad sp credential list/delete` | Manage credentials |
| `az ad sp owner list` | List SP owners |

### Applications

| Command | Description |
|---------|-------------|
| `az ad app create` | Create an application |
| `az ad app delete` | Delete an application |
| `az ad app list` | List applications |
| `az ad app show` | Get application details |
| `az ad app update` | Update an application |
| `az ad app credential reset/list/delete` | Manage app credentials |
| `az ad app permission add/list/grant/admin-consent` | Manage OAuth2 permissions |
| `az ad app federated-credential create/show/list/delete` | Manage federated identity credentials |

### Signed-In User

| Command | Description |
|---------|-------------|
| `az ad signed-in-user show` | Get details of current user |
| `az ad signed-in-user list-owned-objects` | List owned objects |

## `az role` — Role-Based Access Control

### Role Assignments

| Command | Description |
|---------|-------------|
| `az role assignment create` | Create a role assignment |
| `az role assignment list` | List role assignments |
| `az role assignment delete` | Delete role assignments |
| `az role assignment update` | Update a role assignment (via JSON) |
| `az role assignment list-changelogs` | List changelogs in a time window |

Key flags for `az role assignment create`:
- `--role` (required): role name or ID
- `--scope` (required): `/subscriptions/...`, `/subscriptions/.../resourceGroups/...`, or resource-level
- `--assignee`: user sign-in name, SPN, or object ID
- `--assignee-object-id`: bypass Microsoft Graph query
- `--assignee-principal-type`: `ForeignGroup`, `Group`, `ServicePrincipal`, `User`
- `--condition` / `--condition-version` (preview)

```bash
az role assignment create \
  --assignee sp_name \
  --role Reader \
  --scope /subscriptions/.../resourceGroups/MyGroup
```

Key flags for `az role assignment list`:
- `--all`: show all assignments under current subscription (not just subscription-scoped)
- `--include-inherited`: include assignments on parent scopes
- `--include-groups`: include transitive group memberships
- `--fill-principal-name` / `--fill-role-definition-name`: set to `false` for performance (avoids Graph queries)

### Role Definitions

| Command | Description |
|---------|-------------|
| `az role definition create` | Create a custom role definition |
| `az role definition list` | List role definitions |
| `az role definition show` | Show a role definition |
| `az role definition update` | Update a custom role definition |
| `az role definition delete` | Delete a role definition |

Custom role JSON structure:

```json
{
  "Name": "Contoso On-call",
  "Description": "Monitor and restart VMs",
  "Actions": ["Microsoft.Compute/*/read", "Microsoft.Compute/virtualMachines/start/action"],
  "NotActions": [],
  "DataActions": ["Microsoft.Storage/storageAccounts/blobServices/containers/blobs/*"],
  "NotDataActions": ["Microsoft.Storage/storageAccounts/blobServices/containers/blobs/write"],
  "AssignableScopes": ["/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"]
}
```

| Property | Scope | Description |
|----------|-------|-------------|
| `Actions` | Control plane | Allowed management operations (wildcards supported) |
| `NotActions` | Control plane | Denied management operations |
| `DataActions` | Data plane | Allowed data operations |
| `NotDataActions` | Data plane | Denied data operations |
| `AssignableScopes` | Both | Scopes where this role can be assigned |

## `az identity` — Managed Identities

| Command | Description |
|---------|-------------|
| `az identity create` | Create a user-assigned managed identity |
| `az identity list` | List managed identities |
| `az identity show` | Show identity details |
| `az identity delete` | Delete an identity |
| `az identity assign` | Assign identity to a resource |
| `az identity remove` | Remove identity from a resource |
| `az identity federated-credential create/show/list/delete` | Manage federated identity credentials |

## References

- Authentication overview: <https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli>
- Interactive login: <https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli-interactively>
- Service principal: <https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli-service-principal>
- Managed identity: <https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli-managed-identity>
- Subscription management: <https://learn.microsoft.com/en-us/cli/azure/manage-azure-subscriptions-azure-cli>
- RBAC: <https://learn.microsoft.com/en-us/cli/azure/role>
- Entra ID: <https://learn.microsoft.com/en-us/cli/azure/ad>

## Related Skills

- For resource hierarchy (tenant → management group → subscription → resource group), see [resource-hierarchy.md](resource-hierarchy.md)
- For configuration (including `AZURE_CONFIG_DIR` for token cache isolation), see [config.md](config.md)
- For `--subscription` parameter behavior, see [global-parameters.md](global-parameters.md)
