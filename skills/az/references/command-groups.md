# Command Groups

Detailed reference for the major Azure CLI command groups, organized by category.

> **Source of truth**: <https://learn.microsoft.com/en-us/cli/azure/reference-index?view=azure-cli-latest>. For command architecture, see [architecture.md](architecture.md).

## Compute

### `az vm` — Virtual Machines

Manage Linux or Windows virtual machines. Core + Extension. GA.

| Subgroup | Key Commands |
|----------|-------------|
| `vm` | `create`, `list`, `show`, `delete`, `start`, `stop`, `restart`, `deallocate`, `update`, `wait` |
| `vm extension` | `set`, `list`, `show`, `delete`, `wait` |
| `vm image` | `list`, `show`, `list-offers`, `list-publishers`, `list-skus` |
| `vm disk` | `attach`, `detach` |
| `vm nic` | `add`, `remove`, `list`, `show` |
| `vm open-port` | (convenience command for NSG rules) |
| `vm availability-set` | `create`, `list`, `show`, `delete`, `update` |
| `vm run-command` | `invoke`, `list`, `show` |
| `vm secret` | `add`, `list`, `delete` (key vault secrets) |

Key flags for `vm create`:
- `--name`, `--resource-group` (required)
- `--image` (e.g., `Ubuntu2204`, `Win2022Datacenter`)
- `--size` (e.g., `Standard_DS2_v2`)
- `--location`, `--vnet-name`, `--subnet`, `--public-ip-address`, `--nsg`
- `--admin-username`, `--ssh-key-value` (Linux) / `--admin-password` (Windows)
- `--authentication-type` (`ssh`, `password`)
- `--os-disk-size-gb`, `--data-disk-sizes-gb`
- `--tags`, `--no-wait`, `--assign-identity`

Notable: `az vm list` does not show power state by default — use `--show-details` to see it (triggers extra API call per VM).

### `az vmss` — Virtual Machine Scale Sets

| Subgroup | Key Commands |
|----------|-------------|
| `vmss` | `create`, `list`, `show`, `delete`, `update`, `scale`, `start`, `stop`, `restart`, `deallocate` |
| `vmss instance` | `list`, `show`, `delete`, `update` |
| `vmss extension` | `set`, `list`, `show`, `delete` |
| `vmss nic` | `list`, `show` |
| `vmss disk` | `attach`, `detach` |

### `az aks` — Azure Kubernetes Service

| Subgroup | Key Commands |
|----------|-------------|
| `aks` | `create`, `list`, `show`, `delete`, `update`, `scale`, `upgrade`, `stop`, `start` |
| `aks nodepool` | `add`, `list`, `show`, `delete`, `update`, `scale`, `upgrade`, `start`, `stop` |
| `aks addon` | `list`, `show`, `enable`, `disable` |
| `aks approuting` | `enable`, `disable`, `list`, `update`, `show` |
| `aks maintenanceconfiguration` | `add`, `list`, `show`, `delete`, `update` |
| `aks mesh` | `enable`, `disable`, `show` (service mesh) |

Key flags for `aks create`:
- `--name`, `--resource-group` (required)
- `--node-count` (default 3), `--node-vm-size`
- `--enable-managed-identity`, `--attach-acr`
- `--network-plugin` (`azure`, `kubenet`), `--service-cidr`, `--dns-service-ip`
- `--enable-addons` (e.g., `monitoring`, `http_application_routing`)
- `--admin-username`, `--ssh-key-value`
- `--zones` (availability zones)

### `az container` — Container Instances

| Subgroup | Key Commands |
|----------|-------------|
| `container` | `create`, `list`, `show`, `delete`, `restart`, `stop`, `start`, `update`, `exec`, `logs`, `attach` |
| `container app` | (see `containerapp` below) |

### `az containerapp` — Container Apps

| Subgroup | Key Commands |
|----------|-------------|
| `containerapp` | `create`, `list`, `show`, `delete`, `update`, `start`, `stop`, `restart` |
| `containerapp env` | `create`, `list`, `show`, `delete`, `update` (managed environment) |
| `containerapp job` | `create`, `list`, `show`, `delete`, `start`, `stop`, `execute` |
| `containerapp revision` | `list`, `show`, `activate`, `deactivate`, `restart`, `set-mode` |
| `containerapp ingress` | `show`, `enable`, `disable`, `update` |
| `containerapp scale` | `set` |

### `az disk` — Managed Disks

| Subgroup | Key Commands |
|----------|-------------|
| `disk` | `create`, `list`, `show`, `delete`, `update`, `grant-access`, `revoke-access` |
| `disk-access` | `create`, `list`, `show`, `delete`, `update` |
| `disk-encryption-set` | `create`, `list`, `show`, `delete`, `update` |

### `az image` — Custom VM Images

| Subgroup | Key Commands |
|----------|-------------|
| `image` | `create`, `list`, `show`, `delete` |
| `image template` | `create`, `list`, `show`, `delete`, `run`, `cancel`, `wait` (Azure Image Builder) |

### `az sig` — Compute Gallery (formerly Shared Image Gallery)

| Subgroup | Key Commands |
|----------|-------------|
| `sig` | `create`, `list`, `show`, `delete`, `update` |
| `sig image-definition` | `create`, `list`, `show`, `delete`, `update` |
| `sig image-version` | `create`, `list`, `show`, `delete`, `update` |

### `az functionapp` — Function Apps

| Subgroup | Key Commands |
|----------|-------------|
| `functionapp` | `create`, `list`, `show`, `delete`, `start`, `stop`, `restart`, `deploy`, `config` |
| `functionapp deployment` | `source`, `slot`, `user`, `container` |
| `functionapp plan` | `create`, `list`, `show`, `delete`, `update` |
| `functionapp identity` | `assign`, `remove`, `show` |

### `az webapp` — Web Apps

| Subgroup | Key Commands |
|----------|-------------|
| `webapp` | `create`, `list`, `show`, `delete`, `start`, `stop`, `restart`, `deploy`, `config`, `up` |
| `webapp deployment` | `source`, `slot`, `user`, `container` |
| `webapp log` | `tail`, `download`, `config`, `show` |
| `webapp ssl` | `import`, `upload`, `bind`, `unbind`, `delete`, `list`, `show` |
| `webapp identity` | `assign`, `remove`, `show` |
| `webapp vnet-integration` | `add`, `list`, `remove` |

### `az appservice` — App Service Plans

| Subgroup | Key Commands |
|----------|-------------|
| `appservice plan` | `create`, `list`, `show`, `delete`, `update` |
| `appservice vnet-integration` | `add`, `list`, `remove` |

## Networking

### `az network` — Network Resources

The largest and most complex command group. Core + Extension. GA.

| Subgroup | Key Commands |
|----------|-------------|
| `network vnet` | `create`, `list`, `show`, `delete`, `update` |
| `network vnet subnet` | `create`, `list`, `show`, `delete`, `update` |
| `network vnet peering` | `create`, `list`, `show`, `delete`, `update` |
| `network nic` | `create`, `list`, `show`, `delete`, `update` |
| `network nsg` | `create`, `list`, `show`, `delete`, `update` |
| `network nsg rule` | `create`, `list`, `show`, `delete`, `update` |
| `network lb` | `create`, `list`, `show`, `delete`, `update` |
| `network lb rule` | `create`, `list`, `show`, `delete`, `update` |
| `network lb probe` | `create`, `list`, `show`, `delete`, `update` |
| `network lb address-pool` | `create`, `list`, `show`, `delete` |
| `network public-ip` | `create`, `list`, `show`, `delete`, `update` |
| `network dns zone` | `create`, `list`, `show`, `delete`, `update` |
| `network dns record-set` | `a`, `aaaa`, `cname`, `mx`, `txt`, `ptr`, `srv` (each with CRUD) |
| `network application-gateway` | `create`, `list`, `show`, `delete`, `update`, `start`, `stop` |
| `network route-table` | `create`, `list`, `show`, `delete`, `update` |
| `network route-table route` | `create`, `list`, `show`, `delete`, `update` |
| `network private-endpoint` | `create`, `list`, `show`, `delete` |
| `network private-link-service` | `create`, `list`, `show`, `delete`, `update` |
| `network firewall` | `create`, `list`, `show`, `delete`, `update` (extension) |
| `network nat-gateway` | `create`, `list`, `show`, `delete`, `update` |
| `network cross-region-lb` | (cross-region load balancer) |
| `network front-door` | `create`, `list`, `show`, `delete`, `update` (extension) |
| `network traffic-manager` | `profile`, `endpoint` |
| `network watcher` | `configure`, `show`, `list`, `connection-monitor`, `flow-log`, `packet-capture` |

### `az cdn` — Content Delivery Network

| Subgroup | Key Commands |
|----------|-------------|
| `cdn` | `create`, `list`, `show`, `delete`, `update` (profile/endpoint) |
| `cdn custom-domain` | `create`, `list`, `show`, `delete`, `enable-https` |
| `cdn origin` | `create`, `list`, `show`, `delete`, `update` |
| `cdn origin-group` | `create`, `list`, `show`, `delete`, `update` |

### `az signalr` — SignalR Service

| Subgroup | Key Commands |
|----------|-------------|
| `signalr` | `create`, `list`, `show`, `delete`, `update`, `restart` |
| `signalr key` | `list`, `renew` |
| `signalr cors` | `add`, `list`, `remove` |

### `az relay` — Azure Relay Service

| Subgroup | Key Commands |
|----------|-------------|
| `relay namespace` | `create`, `list`, `show`, `delete`, `update` |
| `relay hyco` | `create`, `list`, `show`, `delete`, `update` (hybrid connection) |
| `relay wcfrelay` | `create`, `list`, `show`, `delete` |
| `relay queue` | `create`, `list`, `show`, `delete` |

## Storage & Data

### `az storage` — Azure Storage

| Subgroup | Key Commands |
|----------|-------------|
| `storage account` | `create`, `list`, `show`, `delete`, `update`, `keys list`, `revoke-key` |
| `storage container` | `create`, `list`, `show`, `delete`, `list-blobs` |
| `storage blob` | `upload`, `download`, `list`, `show`, `delete`, `copy` |
| `storage blob upload-batch` | Upload multiple files |
| `storage blob download-batch` | Download multiple blobs |
| `storage share` | `create`, `list`, `show`, `delete` (file share) |
| `storage file` | `upload`, `download`, `list`, `show`, `delete` |
| `storage table` | `create`, `list`, `show`, `delete` |
| `storage queue` | `create`, `list`, `show`, `delete`, `message put`, `message get` |
| `storage cors` | `add`, `list`, `clear` |

Key flags for `storage account create`:
- `--name`, `--resource-group` (required)
- `--sku` (`Standard_LRS`, `Standard_GRS`, `Standard_ZRS`, `Premium_LRS`, etc.)
- `--kind` (`StorageV2`, `BlobStorage`, `FileStorage`, `BlockBlobStorage`)
- `--access-tier` (`Hot`, `Cool`)
- `--enable-hierarchical-namespace` (for Data Lake Gen2)
- `--min-tls-version`, `--allow-blob-public-access`

Notable: `az storage` data-plane commands use `[storage]` config section defaults (`account`, `key`, `sas_token`, `connection_string`).

### `az cosmosdb` — Cosmos DB

| Subgroup | Key Commands |
|----------|-------------|
| `cosmosdb` | `create`, `list`, `show`, `delete`, `update`, `check-name-exists` |
| `cosmosdb database` | `create`, `list`, `show`, `delete` (SQL API) |
| `cosmosdb collection` | `create`, `list`, `show`, `delete`, `update` (SQL API) |
| `cosmosdb sql database` | `create`, `list`, `show`, `delete` |
| `cosmosdb sql container` | `create`, `list`, `show`, `delete`, `update` |
| `cosmosdb mongodb database` | `create`, `list`, `show`, `delete` |
| `cosmosdb mongodb collection` | `create`, `list`, `show`, `delete`, `update` |
| `cosmosdb cassandra keyspace` | `create`, `list`, `show`, `delete` |
| `cosmosdb cassandra table` | `create`, `list`, `show`, `delete` |
| `cosmosdb gremlin database` | `create`, `list`, `show`, `delete` |
| `cosmosdb gremlin graph` | `create`, `list`, `show`, `delete` |
| `cosmosdb table` | `create`, `list`, `show`, `delete` (Table API) |

### `az sql` — Azure SQL

| Subgroup | Key Commands |
|----------|-------------|
| `sql server` | `create`, `list`, `show`, `delete`, `update`, `conn-policy` |
| `sql db` | `create`, `list`, `show`, `delete`, `update`, `copy`, `restore` |
| `sql elastic-pool` | `create`, `list`, `show`, `delete`, `update` |
| `sql firewall-rule` | `create`, `list`, `show`, `delete`, `update` |
| `sql virtual-network-rule` | `create`, `list`, `show`, `delete` |
| `sql mi` | `create`, `list`, `show`, `delete`, `update` (managed instance) |
| `sql midb` | `create`, `list`, `show`, `delete`, `update` (managed instance DB) |
| `sql server ad-admin` | `set`, `list`, `delete` |
| `sql server firewall-rule` | `create`, `list`, `show`, `delete`, `update` |

### `az mysql` — MySQL Database

| Subgroup | Key Commands |
|----------|-------------|
| `mysql server` | `create`, `list`, `show`, `delete`, `update`, `restart` |
| `mysql db` | `create`, `list`, `show`, `delete` |
| `mysql firewall-rule` | `create`, `list`, `show`, `delete`, `update` |
| `mysql flexible-server` | `create`, `list`, `show`, `delete`, `update`, `restart` (flexible server) |
| `mysql flexible-server db` | `create`, `list`, `show`, `delete` |

### `az postgres` — PostgreSQL Database

| Subgroup | Key Commands |
|----------|-------------|
| `postgres server` | `create`, `list`, `show`, `delete`, `update`, `restart` |
| `postgres db` | `create`, `list`, `show`, `delete` |
| `postgres flexible-server` | `create`, `list`, `show`, `delete`, `update`, `restart` |
| `postgres flexible-server db` | `create`, `list`, `show`, `delete` |
| `postgres flexible-server firewall-rule` | `create`, `list`, `show`, `delete` |

### `az redis` — Redis Cache

| Subgroup | Key Commands |
|----------|-------------|
| `redis` | `create`, `list`, `show`, `delete`, `update`, `regenerate-keys`, `force-reboot` |
| `redis firewall-rules` | `create`, `list`, `show`, `delete` |
| `redis patch-schedule` | `set`, `show`, `delete` |

### `az eventhubs` — Event Hubs

| Subgroup | Key Commands |
|----------|-------------|
| `eventhubs namespace` | `create`, `list`, `show`, `delete`, `update` |
| `eventhubs eventhub` | `create`, `list`, `show`, `delete`, `update` |
| `eventhubs consumer-group` | `create`, `list`, `show`, `delete` |
| `eventhubs georecovery-alias` | `set`, `show`, `break-pair`, `exists` |

### `az servicebus` — Service Bus

| Subgroup | Key Commands |
|----------|-------------|
| `servicebus namespace` | `create`, `list`, `show`, `delete`, `update` |
| `servicebus queue` | `create`, `list`, `show`, `delete`, `update` |
| `servicebus topic` | `create`, `list`, `show`, `delete`, `update` |
| `servicebus subscription` | `create`, `list`, `show`, `delete`, `update` |
| `servicebus topic subscription rule` | `create`, `list`, `show`, `delete` |

### `az keyvault` — Key Vault

| Subgroup | Key Commands |
|----------|-------------|
| `keyvault` | `create`, `list`, `show`, `delete`, `update`, `purge`, `recover` |
| `keyvault key` | `create`, `list`, `show`, `delete`, `backup`, `restore`, `encrypt`, `decrypt` |
| `keyvault secret` | `set`, `show`, `list`, `delete`, `backup`, `restore`, `download` |
| `keyvault certificate` | `create`, `list`, `show`, `delete`, `import`, `download` |
| `keyvault storage` | `add`, `list`, `show`, `delete`, `update` |

### `az netappfiles` — Azure NetApp Files

| Subgroup | Key Commands |
|----------|-------------|
| `netappfiles account` | `create`, `list`, `show`, `delete`, `update` |
| `netappfiles pool` | `create`, `list`, `show`, `delete`, `update` |
| `netappfiles volume` | `create`, `list`, `show`, `delete`, `update` |
| `netappfiles snapshot` | `create`, `list`, `show`, `delete` |

## Security & Identity

### `az ad` — Microsoft Entra ID (Azure AD)

| Subgroup | Key Commands |
|----------|-------------|
| `ad user` | `create`, `delete`, `list`, `show`, `update`, `get-member-groups` |
| `ad group` | `create`, `delete`, `list`, `show`, `update`, `member add/list/remove/check`, `owner add/list/remove` |
| `ad sp` | `create`, `create-for-rbac`, `delete`, `list`, `show`, `update`, `credential reset/list/delete` |
| `ad app` | `create`, `delete`, `list`, `show`, `update`, `credential reset/list/delete`, `permission add/list/grant/admin-consent`, `federated-credential create/show/list/delete` |
| `ad signed-in-user` | `show`, `list-owned-objects` |

For authentication details, see [auth.md](auth.md).

### `az role` — Role-Based Access Control

| Subgroup | Key Commands |
|----------|-------------|
| `role assignment` | `create`, `list`, `delete`, `update`, `list-changelogs` |
| `role definition` | `create`, `list`, `show`, `update`, `delete` |

For RBAC details, see [auth.md](auth.md).

### `az identity` — Managed Identities

| Subgroup | Key Commands |
|----------|-------------|
| `identity` | `create`, `list`, `show`, `delete`, `assign`, `remove` |
| `identity federated-credential` | `create`, `list`, `show`, `delete` |

### `az security` — Security Center / Defender for Cloud

| Subgroup | Key Commands |
|----------|-------------|
| `security` | `list-alerts`, `show-alert`, `list-task`, `show-task`, `list-regulatory-compliance-results` |
| `security contact` | `create`, `list`, `show`, `delete` |
| `security auto-provisioning-setting` | `list`, `show`, `update` |
| `security workspace-setting` | `list`, `show`, `update` |

### `az policy` — Azure Policy

| Subgroup | Key Commands |
|----------|-------------|
| `policy assignment` | `create`, `list`, `show`, `delete`, `update` |
| `policy definition` | `create`, `list`, `show`, `delete`, `update` |
| `policy set-definition` | `create`, `list`, `show`, `delete`, `update` |
| `policy exemption` | `create`, `list`, `show`, `delete`, `update` |

### `az keyvault` — See Storage & Data above

### `az lock` — Resource Locks

| Subgroup | Key Commands |
|----------|-------------|
| `lock` | `create`, `list`, `show`, `delete`, `update` |

Lock types: `CanNotDelete`, `ReadOnly`.

## Management & Operations

### `az account` — Subscription Management

| Subgroup | Key Commands |
|----------|-------------|
| `account` | `list`, `show`, `set`, `clear`, `list-locations`, `get-access-token` |
| `account lock` | `create`, `list`, `show`, `delete`, `update` |
| `account management-group` | `create`, `list`, `show`, `delete`, `update`, `subscription add/remove` |
| `account tenant` | `list` (experimental, extension) |

For subscription context details, see [auth.md](auth.md).

### `az group` — Resource Groups

| Subgroup | Key Commands |
|----------|-------------|
| `group` | `create`, `list`, `show`, `delete`, `update`, `exists`, `wait` |
| `group deployment` | `create`, `list`, `show`, `delete`, `cancel`, `validate`, `export-template` |
| `group lock` | `create`, `list`, `show`, `delete`, `update` |

### `az resource` — Generic Resource Management

| Subgroup | Key Commands |
|----------|-------------|
| `resource` | `create`, `list`, `show`, `delete`, `update`, `move`, `tag`, `invoke-action`, `link` |
| `resource lock` | `create`, `list`, `show`, `delete`, `update` |

### `az deployment` — ARM Template Deployment

| Subgroup | Key Commands |
|----------|-------------|
| `deployment` | `create`, `list`, `show`, `delete`, `cancel`, `validate`, `what-if`, `export-template` (subscription scope) |
| `deployment group` | `create`, `list`, `show`, `delete`, `cancel`, `validate`, `what-if` (resource group scope) |
| `deployment mg` | `create`, `list`, `show`, `delete`, `validate` (management group scope) |
| `deployment tenant` | `create`, `list`, `show`, `delete`, `validate` (tenant scope) |
| `deployment-scripts` | `create`, `list`, `show`, `delete` |

### `az stack` — Deployment Stacks

| Subgroup | Key Commands |
|----------|-------------|
| `stack` | `create`, `list`, `show`, `delete`, `export` (subscription scope) |
| `stack group` | `create`, `list`, `show`, `delete`, `export` (resource group scope) |
| `stack mg` | `create`, `list`, `show`, `delete` (management group scope) |
| `stack tenant` | `create`, `list`, `show`, `delete` (tenant scope) |

### `az provider` — Resource Providers

| Subgroup | Key Commands |
|----------|-------------|
| `provider` | `list`, `show`, `register`, `unregister`, `operation list` |
| `provider feature` | `list`, `register`, `show`, `unregister` |

### `az monitor` — Azure Monitor

| Subgroup | Key Commands |
|----------|-------------|
| `monitor log-analytics` | `workspace create/list/show/delete/update`, `query` |
| `monitor log-analytics solution` | `create`, `list`, `show`, `delete` |
| `monitor metrics alert` | `create`, `list`, `show`, `delete`, `update` |
| `monitor diagnostic-settings` | `create`, `list`, `show`, `delete`, `update` |
| `monitor activity-log` | `list`, `alert` |
| `monitor autoscale` | `create`, `list`, `show`, `delete`, `update` |
| `monitor action-group` | `create`, `list`, `show`, `delete`, `update` |

### `az billing` — Billing

| Subgroup | Key Commands |
|----------|-------------|
| `billing account` | `list`, `show`, `update` |
| `billing subscription` | `list`, `show` |
| `billing invoice` | `list`, `show`, `download` |
| `billing period` | `list`, `show` |

### `az tag` — Resource Tags

| Subgroup | Key Commands |
|----------|-------------|
| `tag` | `create`, `list`, `show`, `delete`, `update`, `add-value` |
| `tag list` | List tags on a resource |

### `az feature` — Resource Provider Features

| Subgroup | Key Commands |
|----------|-------------|
| `feature` | `list`, `register`, `show` |

## Containers & DevOps

### `az acr` — Container Registries

| Subgroup | Key Commands |
|----------|-------------|
| `acr` | `create`, `list`, `show`, `delete`, `update`, `check-name` |
| `acr repository` | `list`, `show-tags`, `show-manifests`, `delete`, `update` |
| `acr build` | Build an image (ACR Tasks) |
| `acr task` | `create`, `list`, `show`, `delete`, `update`, `run` |
| `acr pack` | Build using Buildpacks |
| `acr login` | Login to a registry |
| `acr import` | Import an image from another registry |

### `az devops` — Azure DevOps (Extension)

| Subgroup | Key Commands |
|----------|-------------|
| `devops project` | `create`, `list`, `show`, `delete` |
| `devops team` | `create`, `list`, `show`, `delete`, `update` |
| `devops user` | `add`, `list`, `show`, `remove`, `update` |
| `devops pipeline` | (see `pipelines`) |

### `az pipelines` — Azure Pipelines (Extension)

| Subgroup | Key Commands |
|----------|-------------|
| `pipelines` | `create`, `list`, `show`, `delete`, `run`, `update` |
| `pipelines build` | `list`, `show`, `cancel` |
| `pipelines release` | `create`, `list`, `show` |

## AI & Analytics

### `az cognitiveservices` — Cognitive Services

| Subgroup | Key Commands |
|----------|-------------|
| `cognitiveservices account` | `create`, `list`, `show`, `delete`, `update`, `list-keys` |
| `cognitiveservices account deployment` | `create`, `list`, `show`, `delete` |

### `az search` — Azure AI Search

| Subgroup | Key Commands |
|----------|-------------|
| `search service` | `create`, `list`, `show`, `delete`, `update`, `admin-key` |
| `search index` | `create`, `list`, `show`, `delete` (via REST) |
| `search query-key` | `create`, `list`, `delete` |

### `az ml` — Machine Learning (Extension)

| Subgroup | Key Commands |
|----------|-------------|
| `ml workspace` | `create`, `list`, `show`, `delete`, `update` |
| `ml model` | `create`, `list`, `show`, `delete` |
| `ml endpoint` | `create`, `list`, `show`, `delete`, `update`, `invoke` |
| `ml job` | `create`, `list`, `show`, `cancel`, `stream` |
| `ml environment` | `create`, `list`, `show`, `delete` |
| `ml data` | `create`, `list`, `show`, `delete` |
| `ml compute` | `create`, `list`, `show`, `delete`, `start`, `stop`, `restart` |

### `az synapse` — Synapse Analytics

| Subgroup | Key Commands |
|----------|-------------|
| `synapse workspace` | `create`, `list`, `show`, `delete`, `update` |
| `synapse sparkpool` | `create`, `list`, `show`, `delete`, `update` |
| `synapse sqlpool` | `create`, `list`, `show`, `delete`, `pause`, `resume` |
| `synapse notebook` | `create`, `list`, `show`, `delete`, `export` |

### `az databricks` — Databricks (Extension)

| Subgroup | Key Commands |
|----------|-------------|
| `databricks workspace` | `create`, `list`, `show`, `delete`, `update` |
| `databricks access-connector` | `create`, `list`, `show`, `delete`, `update` |

## Integration

### `az eventgrid` — Event Grid

| Subgroup | Key Commands |
|----------|-------------|
| `eventgrid topic` | `create`, `list`, `show`, `delete`, `update` |
| `eventgrid event-subscription` | `create`, `list`, `show`, `delete`, `update` |
| `eventgrid domain` | `create`, `list`, `show`, `delete`, `update` |
| `eventgrid system-topic` | `create`, `list`, `show`, `delete` |

### `az logicapp` — Logic Apps

| Subgroup | Key Commands |
|----------|-------------|
| `logicapp` | `create`, `list`, `show`, `delete`, `start`, `stop`, `restart` |
| `logicapp deployment` | `source`, `slot` |

## Miscellaneous Core Commands

### `az cache` — Deferred Object Cache

| Subgroup | Key Commands |
|----------|-------------|
| `cache` | `list`, `show`, `delete`, `purge` |

Objects cached using `--defer` on create/update commands can be reviewed and deployed in batch using `az cache`.

### `az rest` — Direct REST API

```bash
az rest --method get --url https://management.azure.com/subscriptions/{id}/resourcegroups?api-version=2021-04-01
```

Automatically authenticated. Uses the active subscription's credentials. For details, see [gotchas.md](gotchas.md).

## References

- Full command reference: <https://learn.microsoft.com/en-us/cli/azure/reference-index?view=azure-cli-latest>
- ACR reference: <https://learn.microsoft.com/en-us/cli/azure/acr>
- VM reference: <https://learn.microsoft.com/en-us/cli/azure/vm>
- Network reference: <https://learn.microsoft.com/en-us/cli/azure/network>

## Related Skills

- For command architecture and internals, see [architecture.md](architecture.md)
- For extension-only command groups, see [extensions.md](extensions.md)
- For authentication and RBAC details, see [auth.md](auth.md)
