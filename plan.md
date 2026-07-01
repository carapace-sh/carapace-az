# Implementation Plan: carapace-az

A shell completion provider for the Azure CLI (`az`), following the patterns established by `carapace-aws` and `carapace-gcloud`.

## Architecture Overview

Two binaries, mirroring the established pattern:

| Binary | Purpose |
|--------|---------|
| `carapace-az` | The completion binary — registers az command groups as cobra commands, loads embedded YAML specs lazily, bridges dynamic completions to az's native argcomplete |
| `carapace-spec-az` | The spec generator — converts az CLI command table JSON (from Python introspection) into carapace-spec YAML files |

### Data Flow

```
mcr.microsoft.com/azure-cli Docker image
  → Python introspection script (dump_command_table.py)
    → JSON command table (az_commands.json)
      → carapace-spec-az (Go) converts JSON → YAML specs
        → go generate regenerates az_generated.go (service map)
          → carapace-az embeds YAML at build time
            → runtime: lazy spec loading + argcomplete bridge for dynamic completions
```

## Key Design Decisions

### 1. Scraping Method: Python Introspection (not source patching)

Unlike `carapace-aws` (which clones the repo and parses botocore JSON data files) or `carapace-gcloud` (which uses `gcloud alpha interactive` to dump a JSON command model), az CLI has **no built-in command to dump its command structure**. The Azure CLI is a Python application built on the Knack framework (which extends argparse).

**Approach**: Run a Python script inside the `mcr.microsoft.com/azure-cli` Docker container that uses Azure CLI's internal APIs to load the full command table and serialize it as JSON.

```python
# dump_command_table.py (simplified)
from azure.cli.core import AzCli, MainCommandsLoader
from azure.cli.core.commands import AzCliCommandInvoker
from azure.cli.core.parser import AzCliCommandParser
from azure.cli.core._config import GLOBAL_CONFIG_DIR, ENV_VAR_PREFIX
from azure.cli.core._help import AzCliHelp
from azure.cli.core._output import AzOutputProducer
from azure.cli.core.azlogging import AzCliLogging
from azure.cli.core.file_util import create_invoker_and_load_cmds_and_args
import json

cli = AzCli(cli_name='az', config_dir=GLOBAL_CONFIG_DIR, ...)
create_invoker_and_load_cmds_and_args(cli)

command_table = cli.invocation.commands_loader.command_table
command_group_table = cli.invocation.commands_loader.command_group_table

# Serialize: command name, description, arguments (name, help, required, choices, nargs, type)
# Serialize: command groups (name, help, commands)
output = { ... }
json.dump(output, sys.stdout)
```

**Why not source-patch like Docker?** The Azure CLI is Python, not Go. Carapace's patch-based scraping (`carapace.Gen(cmd)`) only works for Go cobra/urfave/kingpin/kong CLIs. The introspection approach is the standard way to extract az's command structure programmatically.

**Why not parse help text?** Help text scraping is fragile, loses type information, and can't capture argument metadata (required, choices, nargs). Python introspection gives direct access to the argparse argument registry with all metadata.

### 2. Bridge: argcomplete (existing, no new bridge action needed)

The Azure CLI uses **argcomplete** as its completion backend (via Knack's `CLICompletion` class). Carapace-bridge already has `bridge.ActionArgcomplete("az")` which works with az CLI — in fact, the argcomplete bridge's doc comment uses `az` as its example.

**Approach**: Use `bridge.ActionArgcomplete("az")` for dynamic flag value and positional completion, same pattern as carapace-aws uses `bridge.ActionAws` and carapace-gcloud uses `bridge.ActionGcloud`.

No new bridge action (`ActionAz`) needs to be added to carapace-bridge — the existing argcomplete bridge handles it. If we later find az-specific completion needs (e.g., resource group name completion, subscription completion), those can be added as custom Go actions in `pkg/actions/az/`.

### 3. Version Tracking: Docker ecosystem (like carapace-gcloud)

Use a `Dockerfile` + `compose.yaml` with `mcr.microsoft.com/azure-cli:<version>` as the base image. Dependabot tracks the `docker` ecosystem and opens PRs when new az CLI versions are published.

This mirrors carapace-gcloud's approach (which tracks `google/cloud-sdk` Docker image versions). The alternative (carapace-aws's `package.json` npm approach) doesn't apply since az CLI isn't an npm package.

## Directory Structure

```
carapace-az/
├── plan.md                          # This file
├── go.mod
├── go.sum
├── .gitignore
├── AGENTS.md                        # Agent documentation (like carapace-aws/gcloud)
├── LICENSE
├── README.md
├── Dockerfile                       # Scraping image: mcr.microsoft.com/azure-cli + introspection script
├── compose.yaml                     # Docker compose for running the scrape
├── skills/                          # Existing az skill (already present)
│   └── az/
│       ├── SKILL.md
│       └── references/
├── scripts/
│   └── dump_command_table.py        # Python introspection script (run inside Docker)
├── .github/
│   ├── dependabot.yml               # docker + gomod + github-actions ecosystems
│   └── workflows/
│       ├── go.yml                   # Build, test, format, staticcheck
│       └── dependabot.yml           # Spec regeneration on dependabot PRs
├── cmd/
│   ├── carapace-az/                 # The completion binary
│   │   ├── main.go                  # //go:generate directive + entry point
│   │   ├── generate/
│   │   │   └── main.go              # Regenerates az_generated.go from YAML specs
│   │   └── cmd/
│   │       ├── root.go              # Root az command, global flags, service registration
│   │       ├── common/
│   │       │   └── bridge.go        # Delegates to argcomplete bridge
│   │       └── azcli/
│   │           ├── azcli.go         # embed.FS, Services(), Get() — mirrors botocore.go/gcloud.go
│   │           ├── azcli_generated.go  # Auto-generated service map (do not edit)
│   │           ├── az.yaml          # Root spec (global flags, description)
│   │           └── az.<group>.yaml  # Per-command-group YAML specs (~200+ files)
│   └── carapace-spec-az/            # The spec generator binary
│       ├── main.go
│       └── cmd/
│           ├── root.go              # Reads JSON, writes YAML specs
│           └── command.go           # JSON struct definitions + ToSpecCommand() conversion
└── pkg/
    └── actions/
        ├── actions.go               # Registers spec macros
        └── az/
            ├── subscription.go      # Completes subscription names (from az account list)
            └── resourcegroup.go     # Completes resource group names (from az group list)
```

## Implementation Steps

### Phase 1: Project Scaffolding

1. **Initialize Go module**
   - `go mod init github.com/carapace-sh/carapace-az`
   - Go version: `1.25.0` (matching carapace-aws)
   - Dependencies: `carapace`, `carapace-bridge`, `carapace-spec`, `cobra`, `pflag` (replaced), `yaml.v3`, `sentences`
   - Replace `spf13/pflag` with `carapace-sh/carapace-pflag`

2. **Create base files**: `.gitignore`, `LICENSE`, `README.md`, `AGENTS.md`

3. **Create directory structure** as outlined above

### Phase 2: Scraping Pipeline

4. **Write `scripts/dump_command_table.py`**
   - Import Azure CLI core classes
   - Create `AzCli` instance and load full command table via `create_invoker_and_load_cmds_and_args()`
   - Iterate `command_table` (dict of `command_name → AzCliCommand`)
   - For each command, extract:
     - Command name (e.g., `vm create`, `network vnet subnet list`)
     - Description/help text
     - Arguments: name, options_list (`--name`, `-n`), help, required, choices, type, nargs
   - Iterate `command_group_table` (dict of group names with help/commands)
   - Serialize as JSON to stdout
   - Handle extensions (load extension command loaders too)

5. **Write `Dockerfile`**
   ```dockerfile
   ARG VERSION=2.71.0
   FROM mcr.microsoft.com/azure-cli:${VERSION}
   ADD scripts/dump_command_table.py /
   CMD ["python", "/dump_command_table.py"]
   ```

6. **Write `compose.yaml`**
   ```yaml
   services:
     az:
       build:
         context: .
         args:
           VERSION: 2.71.0
   ```

7. **Test the scrape locally**
   - `docker compose build && docker compose run --rm az > az_commands.json`
   - Verify JSON structure is complete and parseable

### Phase 3: Spec Generator (`carapace-spec-az`)

8. **Write `cmd/carapace-spec-az/cmd/command.go`**
   - Define Go structs matching the JSON output from `dump_command_table.py`:
     - `Cli` — top-level metadata (az version)
     - `Command` — single command with name, description, arguments, subcommands
     - `Arg` — argument with name, help, required, choices, type, nargs
   - Implement `ToSpecCommand()` that converts to `carapace-spec` `command.Command`
   - Use `sentences` tokenizer for first-sentence extraction (like carapace-gcloud)
   - Convert az argument names to kebab-case flag names (az uses `--resource-group` style)

9. **Write `cmd/carapace-spec-az/cmd/root.go`**
   - Accept JSON file path as positional argument
   - Flags: `--target` (output directory), `--stdout`, `--no-doc` (mutually exclusive)
   - Parse JSON → build command tree → group commands by top-level group
   - Write one YAML file per top-level command group: `az.<group>.yaml`
   - Write root spec: `az.yaml` (global flags, description)
   - Add `# yaml-language-server: $schema=https://carapace.sh/schemas/command.json` header

10. **Write `cmd/carapace-spec-az/main.go`**
    - Simple entry point calling `cmd.Execute()`

### Phase 4: Completion Binary (`carapace-az`)

11. **Write `cmd/carapace-az/cmd/azcli/azcli.go`**
    - `//go:embed *.yaml` for all spec files
    - `Services()` returns the service map
    - `Get(name)` reads and unmarshals a specific YAML spec
    - Direct mirror of `botocore.go` / `gcloud.go`

12. **Write `cmd/carapace-az/cmd/common/bridge.go`**
    - `ActionBridgeAzCompleter()` delegates to `bridge.ActionArgcomplete("az")`
    - Same arg passthrough pattern as carapace-aws/gcloud

13. **Write `cmd/carapace-az/cmd/root.go`**
    - Define root `az` command with global persistent flags:
      - `--output` / `-o` (json, jsonc, table, tsv, yaml, yamlc, none)
      - `--query` (JMESPath query string)
      - `--subscription` (subscription name or ID)
      - `--verbose`
      - `--debug`
      - `--only-show-errors`
      - `--help` / `-h`
      - `--version`
    - Flag completion for `--output` (static values), `--subscription` (custom action)
    - Iterate `azcli.Services()` to create per-group cobra commands
    - Each group command uses `PreRun` to load its YAML spec and register subcommands via `spec.Command.ToCobra()`
    - `PreInvoke` hook: delegate non-bool flag completion to argcomplete bridge when no static completion is defined
    - Top-level non-group commands (`login`, `logout`, `rest`, `version`, `configure`, `interactive`, `find`, `feedback`, `upgrade`, `survey`, `init`) with `DisableFlagParsing` and full bridge delegation
    - `spec.Register(rootCmd)` at the end

14. **Write `cmd/carapace-az/generate/main.go`**
    - Scan `az.*.yaml` files in `../cmd/azcli/`
    - Extract name + description from each
    - Write `azcli_generated.go` with `init()` populating `services` map
    - Run `go fmt` on the generated file
    - Direct mirror of carapace-aws/gcloud generate scripts

15. **Write `cmd/carapace-az/main.go`**
    - `//go:generate sh -c "go run -C ./generate ."`
    - Call `cmd.Execute()`

### Phase 5: Custom Actions

16. **Write `pkg/actions/az/subscription.go`**
    - `ActionSubscriptions()` — completes subscription names from `az account list --output json`
    - Parse JSON output, extract `name` and `id`, offer both with descriptions

17. **Write `pkg/actions/az/resourcegroup.go`**
    - `ActionResourceGroups()` — completes resource group names from `az group list --output json`
    - Parse JSON output, extract `name` and `location`

18. **Write `pkg/actions/actions.go`**
    - Register `Subscriptions` and `ResourceGroups` as spec macros

### Phase 6: CI/CD & Dependabot

19. **Write `.github/dependabot.yml`**
    ```yaml
    version: 2
    updates:
      - package-ecosystem: "gomod"
        directory: "/"
        schedule:
          interval: "daily"
      - package-ecosystem: "github-actions"
        directory: "/"
        schedule:
          interval: "daily"
      - package-ecosystem: "docker"
        directory: "/"
        schedule:
          interval: "daily"
    ```

20. **Write `.github/workflows/dependabot.yml`**
    - Trigger: `pull_request` by `dependabot[bot]`
    - `az` job (in `ghcr.io/carapace-sh/go:1.25.4` container):
      1. Deep clone
      2. `docker compose build && docker compose run --rm az > az_commands.json`
      3. `go run -C cmd/carapace-spec-az . --target cmd/carapace-az/cmd/azcli --no-doc az_commands.json`
      4. `go generate ./cmd/carapace-az`
      5. Commit and push updated specs to PR branch
    - `auto-merge` job: approve and auto-merge

21. **Write `.github/workflows/go.yml`**
    - Trigger: `pull_request`, `push`
    - Build all `cmd/` dirs, test, gofmt check, staticcheck, coverage
    - GoReleaser on tag pushes (if desired)

### Phase 7: Testing & Validation

22. **Write `cmd/carapace-az/main_test.go`**
    - Integration tests: build `carapace-az`, run completion for major command groups, diff against `az` CLI's own completion output
    - Pattern from carapace-aws: test binary that runs per-service completion comparison

23. **Manual validation**
    - Build the completer: `go build -C cmd/carapace-az .`
    - Test: `carapace-az _carapace spec` should output the full spec
    - Test shell completion integration

## Key Considerations

### az CLI Command Structure (Hierarchical)

Unlike AWS (flat `aws <service> <operation>`) or gcloud (`gcloud <service> <subgroup> <operation>`), az CLI has deeply nested command groups:

```
az network vnet subnet create
az network vnet subnet list
az network nic create
az vm create
az storage account create
```

The spec generator must handle this nesting properly — either:
- **Option A**: One YAML per top-level group (`az.network.yaml`), with nested subcommands in the spec (like carapace-gcloud)
- **Option B**: One YAML per leaf command group (more granular files)

**Recommendation**: Option A (one file per top-level group), matching carapace-gcloud's approach. This keeps the file count manageable (~50-80 files vs 500+ for per-leaf-group) and lets the spec's nested `commands` structure handle the hierarchy.

### Extensions

az CLI has a dynamic extension system. Extensions add new command groups (e.g., `az devops`, `az kusto`, `az ml`). For the initial implementation:
- Scrape only **core** command modules (built-in)
- Extensions can be added later via customizations or by installing extensions in the Dockerfile before scraping

### Global Parameters

az CLI global parameters must be registered as persistent flags on the root command, with static completion where applicable:

| Flag | Completion |
|------|-----------|
| `--output` / `-o` | Static: `json`, `jsonc`, `table`, `tsv`, `yaml`, `yamlc`, `none` |
| `--query` | None (free-form JMESPath) |
| `--subscription` | Dynamic: `ActionSubscriptions()` |
| `--verbose` | Boolean |
| `--debug` | Boolean |
| `--only-show-errors` | Boolean |

### Argument Name Conventions

az CLI uses `--kebab-case` flag names consistently (e.g., `--resource-group`, `--vnet-name`). The Python introspection gives us the `options_list` directly (e.g., `['--name', '-n']`), so no CamelCase-to-kebab conversion is needed (unlike carapace-aws which needs extensive CamelCaseToDash fixes).

### Customizations

Like carapace-aws has `cmd/carapace-spec-botocore/cmd/customizations/`, we may need a customizations directory for:
- Commands that aren't in the command table (e.g., `az interactive`, `az find`)
- Flag overrides for specific commands
- Removing deprecated commands

Start with a minimal customizations framework and add as needed.

## Dependency Summary

| Dependency | Version (approx) | Purpose |
|-----------|-----------------|---------|
| `carapace` | `v1.13.0` | Shell completion engine |
| `carapace-bridge` | `v1.6.1` | argcomplete bridge for az |
| `carapace-spec` | `v1.7.1` | YAML spec types, ToCobra() conversion |
| `cobra` | `v1.10.2` | CLI framework |
| `carapace-pflag` | `v1.1.0` | Patched pflag (replace directive) |
| `yaml.v3` | `v3.0.1` | YAML marshaling |
| `sentences` | `v1.1.2` | Sentence tokenizer for descriptions |
| `mcr.microsoft.com/azure-cli` | `2.71.0+` | Docker image for scraping |

## Open Questions

1. **Extension support**: Should we install common extensions (e.g., `azdev`, `ai-examples`, `ssh`) in the Dockerfile to scrape their commands too? Or keep to core only for v1?
   - **Recommendation**: Core only for v1. Add extensions later.

2. **Integration tests**: carapace-aws has extensive per-service integration tests using `amazon/aws-cli` containers. Should we do the same with `mcr.microsoft.com/azure-cli`?
   - **Recommendation**: Start with basic build tests. Add per-group integration tests after v1 is stable.

3. **GoReleaser**: Should we set up GoReleaser for binary releases (like carapace-aws) or skip it for now?
   - **Recommendation**: Set up `.goreleaser.yml` but don't block v1 on it.

4. **Bridge action**: Should we add a dedicated `ActionAz` to carapace-bridge (like `ActionAws`, `ActionGcloud`) or keep using `ActionArgcomplete("az")`?
   - **Recommendation**: Use `ActionArgcomplete("az")` for now. If az-specific completion handling is needed later, add a dedicated action.
