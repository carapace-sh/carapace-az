---
name: az
description: >
  Use when working with the Azure CLI (az) — command structure, global parameters,
  configuration system, authentication, resource hierarchy, output formatting, extensions,
  and command groups. Triggers on: "az", "azure cli", "az cli", "az login", "az account",
  "az vm", "az network", "az storage", "az aks", "az sql", "az keyvault", "az group",
  "az resource", "az ad", "az role", "az acr", "az appservice", "az webapp", "az functionapp",
  "az cosmosdb", "az containerapp", "az deployment", "az config", "az extension",
  "--query", "JMESPath", "az rest", "az interactive", "service principal", "managed identity",
  "resource group", "subscription", "AZURE_CONFIG_DIR", "AZURE_", "knack", "az extension add".
user-invocable: true
---

# Azure CLI (az) In-Depth Reference

Comprehensive reference for the Azure CLI (`az`) — the cross-platform command-line tool for managing Azure resources and services.

## Data Flow

```
az command line
  → argument parsing (global parameters + command parameters)
    → configuration resolution (CLI params > env vars > config file)
      → credential lookup (user / service principal / managed identity / WAM)
        → subscription context resolution (--subscription or active subscription)
          → API request (Azure Resource Manager or data-plane)
            → response processing (--query JMESPath)
              → output formatting (--output)
                → stdout (results) / stderr (warnings, errors)
```

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| command structure, command pattern, command groups, dispatch, CLI architecture, knack, command loader, command table, argument registry, parser hierarchy, module system, interactive mode, az interactive | [references/architecture.md](references/architecture.md) |
| service groups, command groups, vm, network, storage, aks, sql, keyvault, acr, appservice, webapp, functionapp, cosmosdb, containerapp, monitor, policy, backup, redis, eventhubs, servicebus, postgres, mysql, cdn, signalr, search | [references/command-groups.md](references/command-groups.md) |
| global parameters, --output, --query, --subscription, --debug, --verbose, --only-show-errors, -o, -h, parameter precedence | [references/global-parameters.md](references/global-parameters.md) |
| config, configuration, az config, az configure, az init, AZURE_CONFIG_DIR, environment variables, AZURE_, defaults, core.output, extension.use_dynamic_install, cloud.name, logging | [references/config.md](references/config.md) |
| auth, authentication, credentials, az login, az logout, service principal, managed identity, WAM, MSAL, az account, az account set, az account list, az account get-access-token, subscription context, multi-tenant, az ad, Entra ID, az role, RBAC | [references/auth.md](references/auth.md) |
| resource hierarchy, tenant, management group, subscription, resource group, resource ID, az group, az account management-group, resource provider, az provider, ARM, az deployment, az stack | [references/resource-hierarchy.md](references/resource-hierarchy.md) |
| output formatting, --output, --query, JMESPath, json, jsonc, table, tsv, yaml, yamlc, none, filter expressions, projections, multiselect, pipe expressions, JMESPath functions | [references/output-formatting.md](references/output-formatting.md) |
| extensions, az extension, az extension add, az extension list, dynamic install, azext_, wheel, extension index, extension compatibility, core vs extension, preview extensions | [references/extensions.md](references/extensions.md) |
| gotchas, edge cases, known issues, pitfalls, shell quoting, TSV ordering, MSAL token cache, extension auto-update, --no-wait, az rest, az vm list --show-details, telemetry, INI format, deprecated commands | [references/gotchas.md](references/gotchas.md) |

## Quick Guide

- **How does the az command structure work?** → [references/architecture.md](references/architecture.md)
- **What service command groups are available?** → [references/command-groups.md](references/command-groups.md)
- **What global parameters exist and how do they work?** → [references/global-parameters.md](references/global-parameters.md)
- **How do I configure az CLI?** → [references/config.md](references/config.md)
- **How do I authenticate with az CLI?** → [references/auth.md](references/auth.md)
- **How does the Azure resource hierarchy work?** → [references/resource-hierarchy.md](references/resource-hierarchy.md)
- **How do I format and query output?** → [references/output-formatting.md](references/output-formatting.md)
- **How do extensions work?** → [references/extensions.md](references/extensions.md)
- **What are common gotchas and pitfalls?** → [references/gotchas.md](references/gotchas.md)
- **How do I use JMESPath queries?** → [references/output-formatting.md](references/output-formatting.md)
- **How do I authenticate with a service principal?** → [references/auth.md](references/auth.md)
- **How do I manage subscriptions?** → [references/auth.md](references/auth.md)
- **How do I install and manage extensions?** → [references/extensions.md](references/extensions.md)
- **How do I use az rest for API calls?** → [references/gotchas.md](references/gotchas.md)

## Cross-Project References

- For shell completion integration with az (carapace-bridge, carapace-spec), see the **carapace** skill and **carapace-dev** skill.
- For cobra command structure patterns used by carapace-az, see the **cobra** skill.
- For YAML spec format used by carapace-spec-az, see the **carapace** skill → spec documentation.
- For JMESPath specification details, see <https://jmespath.org/specification.html>.
