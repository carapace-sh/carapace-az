# Resource Hierarchy

The Azure resource hierarchy — tenants, management groups, subscriptions, resource groups, and resources — and the commands that manage them.

> **Source of truth**: <https://learn.microsoft.com/en-us/azure/governance/management-groups/overview> and <https://learn.microsoft.com/en-us/cli/azure/group>. For authentication and subscription context, see [auth.md](auth.md).

## Hierarchy Structure

```
Tenant (Microsoft Entra ID instance)
  └── Management Groups (optional grouping of subscriptions)
       └── Subscriptions
            └── Resource Groups
                 └── Resources (VMs, storage accounts, databases, etc.)
```

Policy and access (RBAC) inheritance flows **downward**: policies and role assignments at a parent scope apply to all child scopes.

## Tenant

A tenant is an instance of Microsoft Entra ID representing an organization.

| Aspect | Details |
|--------|---------|
| Relationship | Root of the hierarchy; contains management groups and subscriptions |
| Identification | Tenant ID (GUID) or verified domain name (e.g., `mycompany.onmicrosoft.com`) |
| Get current tenant | `az account show --query tenantId` |
| List accessible tenants | `az account tenant list` (experimental, extension) |
| Sign in to a specific tenant | `az login --tenant <tenant-id>` |

## Management Groups

Containers for managing access, policies, and compliance across multiple subscriptions.

| Command | Description |
|---------|-------------|
| `az account management-group create` | Create a management group |
| `az account management-group list` | List management groups |
| `az account management-group show` | Show details |
| `az account management-group update` | Update |
| `az account management-group delete` | Delete |
| `az account management-group subscription add` | Add subscription to MG |
| `az account management-group subscription remove` | Remove subscription from MG |
| `az account management-group subscription list` | List subscriptions in MG |

```bash
# Create a management group
az account management-group create --name "prod-mg" --display-name "Production"

# Add a subscription to it
az account management-group subscription add --name "prod-mg" --subscription "MyProdSub"

# List hierarchy
az account management-group list --output table
```

Management groups can be **nested** up to 6 levels deep. The root management group is at level 0.

## Subscriptions

Agreements with Microsoft to use Azure services. Every resource is associated with exactly one subscription.

| Command | Description |
|---------|-------------|
| `az account list` | List all accessible subscriptions |
| `az account show` | Show active subscription |
| `az account set` | Set active subscription |
| `az account list-locations` | List supported regions |
| `az account clear` | Clear all subscriptions from cache |

```bash
# List subscriptions
az account list --output table

# Switch active subscription
az account set --subscription "My Dev Subscription"

# Switch by ID
az account set --subscription "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

# Get active subscription ID for scripting
subId=$(az account show --query id --output tsv)
```

### Subscription Context

- **Active subscription**: All CLI commands run against the active subscription set by `az account set`
- **Per-command override**: `--subscription` flag on any command targets a different subscription without changing the default
- **Cross-tenant**: Changing to a subscription in a different tenant changes the active tenant context
- **Caching**: Subscription info is cached locally; if new permissions are granted, you may need to `az logout && az login` or `az account clear && az login`

## Resource Groups

Logical containers for related resources. Every resource must belong to exactly one resource group.

| Command | Description |
|---------|-------------|
| `az group create` | Create a resource group |
| `az group list` | List resource groups |
| `az group show` | Show details |
| `az group delete` | Delete a resource group (and all its resources) |
| `az group update` | Update (e.g., add tags) |
| `az group exists` | Check if a resource group exists |
| `az group wait` | Block until a condition is met |

```bash
# Create a resource group
az group create --name MyResourceGroup --location eastus

# List resource groups
az group list --output table

# Set as default
az config set defaults.group=MyResourceGroup

# Delete a resource group (destructive — deletes all contained resources)
az group delete --name MyResourceGroup --no-wait
```

### Resource Group Scope

Resource groups are scoped to a single subscription and a single region. The region specified at creation determines where the resource group **metadata** is stored, not where the resources themselves are stored (each resource can be in a different region).

## Resources

Individual Azure services managed by their respective command groups (`az vm`, `az storage`, etc.). The generic `az resource` commands work with any resource type.

| Command | Description |
|---------|-------------|
| `az resource list` | List resources (optionally by resource group or tag) |
| `az resource show` | Show resource details |
| `az resource create` | Create a resource (generic) |
| `az resource update` | Update a resource |
| `az resource delete` | Delete a resource |
| `az resource move` | Move resources between resource groups |
| `az resource tag` | Apply tags to a resource |
| `az resource invoke-action` | Invoke a custom action on a resource |

```bash
# List all resources in a resource group
az resource list --resource-group MyResourceGroup --output table

# List by tag
az resource list --tag environment=prod --output table

# Get a specific resource by ID
az resource show --ids /subscriptions/.../resourceGroups/MyGroup/providers/Microsoft.Compute/virtualMachines/MyVm
```

## Resource IDs

All Azure resources have a unique resource ID following this pattern:

```
/subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/{providerNamespace}/{resourceType}/{resourceName}
```

Examples:
```
/subscriptions/abc-123/resourceGroups/MyGroup/providers/Microsoft.Compute/virtualMachines/MyVm
/subscriptions/abc-123/resourceGroups/MyGroup/providers/Microsoft.Storage/storageAccounts/mystorage
/subscriptions/abc-123/resourceGroups/MyGroup/providers/Microsoft.Network/virtualNetworks/MyVnet
```

Child resources have extended paths:
```
/subscriptions/{id}/resourceGroups/{rg}/providers/Microsoft.Network/virtualNetworks/{vnet}/subnets/{subnet}
```

Resource IDs are used by:
- `--ids` parameter on many commands
- `--scope` on `az role assignment create`
- `--source-resource` on `az resource move`

## Resource Providers

Resource providers are services that offer Azure resource types (e.g., `Microsoft.Compute`, `Microsoft.Network`, `Microsoft.Storage`).

| Command | Description |
|---------|-------------|
| `az provider list` | List all resource providers |
| `az provider show --namespace Microsoft.Compute` | Show provider details |
| `az provider register --namespace Microsoft.PolicyInsights` | Register a provider |
| `az provider unregister --namespace Microsoft.PolicyInsights` | Unregister a provider |
| `az provider operation list` | List operations for a provider |

```bash
# List registered providers
az provider list --query "[?registrationState=='Registered'].namespace" --output table

# Register a new provider
az provider register --namespace Microsoft.ContainerInstance
```

## ARM Template Deployment

Deploy resources using ARM templates (JSON) or Bicep files.

| Scope | Command Group |
|-------|--------------|
| Subscription | `az deployment create` |
| Resource Group | `az deployment group create` |
| Management Group | `az deployment mg create` |
| Tenant | `az deployment tenant create` |

```bash
# Deploy a Bicep file to a resource group
az deployment group create \
  --resource-group MyResourceGroup \
  --template-file main.bicep \
  --parameters env=prod

# Deploy an ARM template from a URL
az deployment group create \
  --resource-group MyResourceGroup \
  --template-uri https://example.com/template.json

# What-if (preview changes before deploying)
az deployment group what-if \
  --resource-group MyResourceGroup \
  --template-file main.bicep
```

### `az stack` — Deployment Stacks

Deployment stacks provide atomic deployment and deletion of resource collections:

```bash
az stack group create --name MyStack --resource-group MyGroup --template-file main.bicep
az stack group list --output table
az stack group show --name MyStack --resource-group MyGroup
az stack group delete --name MyStack --resource-group MyGroup
```

### `az ts` — Template Specs

Store and manage ARM templates as versioned template specs:

```bash
az ts create --name MyTemplateSpec --resource-group MyGroup --template-file main.json --version "1.0"
az ts list --output table
az ts show --name MyTemplateSpec --resource-group MyGroup --version "1.0"
```

## Policy and Access Inheritance

RBAC role assignments and Azure Policy assignments flow **downward** through the hierarchy:

```
Tenant (role assignments, policy assignments)
  ↓ applies to
Management Group (role assignments, policy assignments)
  ↓ applies to
Subscription (role assignments, policy assignments)
  ↓ applies to
Resource Group (role assignments, policy assignments)
  ↓ applies to
Resource (role assignments, policy assignments)
```

Effective permissions = union of all role assignments in the ancestry chain. For RBAC details, see [auth.md](auth.md). For Azure Policy commands, see [command-groups.md](command-groups.md).

## References

- Management groups: <https://learn.microsoft.com/en-us/azure/governance/management-groups/overview>
- Resource groups: <https://learn.microsoft.com/en-us/cli/azure/group>
- ARM templates: <https://learn.microsoft.com/en-us/azure/azure-resource-manager/templates/deploy-cli>
- Deployment stacks: <https://learn.microsoft.com/en-us/azure/azure-resource-manager/bicep/deployment-stacks>
- Resource IDs: <https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/resource-id-templates>

## Related Skills

- For authentication and RBAC role assignments, see [auth.md](auth.md)
- For subscription management commands, see [auth.md](auth.md)
- For command group details, see [command-groups.md](command-groups.md)
