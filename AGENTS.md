# AGENTS.md

## Project Overview

`carapace-az` is a shell completion provider for the Azure CLI (`az`). It enriches the az completer by combining static YAML command specs (embedded at build time) with dynamic completions bridged from az's native argcomplete completer at runtime.

## Architecture

There are two separate binaries:

### `carapace-az` (the completer)

- **Entry point**: `cmd/carapace-az/main.go`
- **Root command**: `cmd/carapace-az/cmd/root.go` — defines the `az` root command with all global persistent flags (matching az CLI's own), then dynamically registers one cobra subcommand per command group (e.g. `vm`, `network`, `storage`) from the `azcli.Services()` map.
- **YAML specs**: `cmd/carapace-az/cmd/azcli/az.*.yaml` — one file per top-level command group, embedded via `//go:embed *.yaml`. Each spec is a `carapace-spec` `Command` YAML describing nested subcommands, flags, and static flag completions.
- **Spec loading**: `cmd/carapace-az/cmd/azcli/azcli.go` — `Services()` returns the group→description map; `Get(name)` reads and unmarshals a specific YAML spec.
- **Generated file**: `cmd/carapace-az/cmd/azcli/azcli_generated.go` — auto-generated `init()` populating the `services` map. **Do not edit by hand**; regenerate with `go generate` (see below).
- **Bridge**: `cmd/carapace-az/cmd/common/bridge.go` — `ActionBridgeAzCompleter()` delegates flag and positional completion to az's native completer via `carapace-bridge`'s argcomplete bridge. The `PreInvoke` hook in `root.go` replaces non-bool flag actions with this bridge action when no static completion is defined.
- **Top-level commands**: Non-group commands (`login`, `logout`, `rest`, `version`, `configure`, `interactive`, `find`, `feedback`, `upgrade`, `survey`, `init`) bypass flag parsing and delegate entirely to the bridge completer.
- **Usage suppression**: `rootCmd.SetUsageFunc` returns nil to suppress cobra usage output (this is a completer, not a real CLI).

### `carapace-spec-az` (the spec generator)

- **Entry point**: `cmd/carapace-spec-az/main.go`
- **Root command**: `cmd/carapace-spec-az/cmd/root.go` — reads a JSON file (az CLI command table from Python introspection), converts it to `carapace-spec` YAML, and writes one file per top-level command group.
- **JSON types**: `cmd/carapace-spec-az/cmd/command.go` — defines `CliData`, `CommandData`, `ArgumentData`, `GroupData` structs matching the JSON schema from `dump_command_table.py`.
- **Conversion logic**: `buildSpec()` groups commands by top-level group, `insertCommand()` builds the nested command tree from flat space-separated command names, `convertArgument()` maps az arguments to carapace-spec flags. Uses `sentences` tokenizer to extract first sentence of descriptions.
- **Flags**: `--target` (output directory), `--stdout` (print to stdout), `--no-doc` (strip documentation). `--target` and `--stdout` are mutually exclusive.

### Data flow

1. **Command table extraction**: `dump_command_table.py` runs inside `mcr.microsoft.com/azure-cli` Docker container, uses knack's `command_table` API via `get_default_cli()` + `create_invoker_and_load_cmds_and_args()` to dump all commands, arguments, and groups (including `long-summary` help text) as JSON
2. **Spec generation (completion)**: `carapace-spec-az --no-doc` converts JSON to YAML specs in `cmd/carapace-az/cmd/azcli/` (no documentation fields, keeping specs small)
3. **Spec generation (man pages)**: `carapace-spec-az --stdout` emits a full spec with `documentation.command` and `documentation.flag` populated from knack's `long-summary`; `carapace-man update` splits it into per-subcommand files in `man/cmd/az/`
4. **Code generation**: `go generate` scans YAML files and produces `azcli_generated.go` with the service map
5. **Runtime**: `carapace-az` embeds YAML specs, registers cobra commands from the service map, loads specs lazily in `PreRun`, converts them to cobra commands via `spec.Command.ToCobra()`, and bridges dynamic completions to az's argcomplete completer

## Directory Structure

```
cmd/
  carapace-az/               # Main completion binary
    cmd/
      root.go                # CLI structure: root + group + operation commands
      azcli/                 # YAML spec files (az.<group>.yaml)
        azcli.go             # Loads embedded YAML files
        azcli_generated.go   # Group name -> description map (generated)
        az.yaml              # Root spec
        az.<group>.yaml      # Per-group command specs
      common/
        bridge.go            # Delegates to carapace-bridge argcomplete
    generate/
      main.go               # Regenerates azcli_generated.go from YAML specs
  carapace-spec-az/          # Spec generator binary
    cmd/
      root.go               # Parses JSON -> YAML specs
      command.go            # JSON struct definitions
scripts/
  dump_command_table.py      # Python introspection script (run in Docker)
pkg/actions/az/              # Go completion actions (subscription, resourcegroup)
man/cmd/az/                  # Man page specs (documentation.command/flag from knack long-summary)
Dockerfile                   # Scraping image: mcr.microsoft.com/azure-cli
compose.yaml                 # Docker compose for running the scrape
```

## Commands

```sh
# Build the completer
go build -o cmd/carapace-az/carapace-az ./cmd/carapace-az

# Build the spec generator
go build -o cmd/carapace-spec-az/carapace-spec-az ./cmd/carapace-spec-az

# Regenerate azcli_generated.go from YAML specs (run from repo root)
go generate ./cmd/carapace-az

# Run tests
go test ./...

# Docker: build image that runs dump_command_table.py to produce the JSON command table
docker compose build
docker compose run --rm az > az_commands.json

# Full spec regeneration pipeline
docker compose build && docker compose run --rm az > az_commands.json
go run -C cmd/carapace-spec-az . --target cmd/carapace-az/cmd/azcli --no-doc az_commands.json
go run -C cmd/carapace-spec-az . --stdout az_commands.json > /tmp/az-full-spec.yaml
carapace-man update /tmp/az-full-spec.yaml man/cmd/az
go generate ./cmd/carapace-az
```

## The `go generate` pipeline

`main.go` has `//go:generate sh -c "go run -C ./generate ."`. This runs `cmd/carapace-az/generate/main.go`, which:
1. Reads all `az.*.yaml` files in `cmd/azcli/`
2. Extracts `name` and `description` from each
3. Writes `azcli_generated.go` with an `init()` that populates the `services` map
4. Runs `go fmt` on the generated file

Note: `go generate` only regenerates `azcli_generated.go` from existing YAML specs. The YAML specs themselves are generated by `carapace-spec-az` from the Docker-scraped JSON.

## Key Dependencies

- **`carapace`** — shell completion engine; provides `carapace.Gen()`, `ActionCallback`, `ActionMap`, etc.
- **`carapace-spec`** — YAML-driven command specs; `spec.Command` unmarshals YAML and converts to cobra commands via `ToCobra()`. Also provides `spec.Register(rootCmd)` for spec-based invocation. `spec.AddMacro()` registers custom completion macros.
- **`carapace-bridge`** — bridges completions from other completers; `bridge.ActionArgcomplete("az")` invokes az's own argcomplete-based completion.
- **`pflag` (replaced)** — uses `carapace-sh/carapace-pflag` (a fork) via `replace` directive in `go.mod`.
- **`sentences`** — sentence tokenizer for extracting first sentence of flag/command descriptions.

## Custom Actions

Custom Go completion actions live in `pkg/actions/az/`:

| Action | Macro | Description |
|--------|-------|-------------|
| `ActionSubscriptions()` | `$az.subscriptions` | Completes subscription names from `az account list` |
| `ActionResourceGroups()` | `$az.resourcegroups` | Completes resource group names from `az group list` |

Macros are registered in `pkg/actions/actions.go` via `spec.AddMacro()` and can be referenced in YAML specs.

## Spec Format

AZ command specs are YAML files following this schema:

```yaml
# yaml-language-server: $schema=https://carapace.sh/schemas/command.json
name: vm
description: Manage virtual machines.
commands:
  - name: create
    description: Create a virtual machine.
    flags:
      -g, --resource-group=!: Name of resource group.
      -n, --name=!: Name of the virtual machine.
      --image=!: The name of the operating system image.
      --no-wait: Do not wait for the long-running operation to finish.
```

Flag conventions:
- `--flag=!` = required flag
- `--flag=` = optional flag (takes argument)
- `--flag` = boolean flag (no argument)
- `--flag*` = repeatable flag
- `nargs: -1` = variadic (multiple values allowed)
- `-x, --long=` = shorthand + longhand

## Scraping Pipeline

### Python Introspection (`scripts/dump_command_table.py`)

Uses knack's public `command_table` API (not raw argparse introspection):

1. `get_default_cli()` creates the AzCli instance (stable across az CLI versions)
2. `create_invoker_and_load_cmds_and_args(cli)` loads the full command table
3. Iterates `command_table` — each entry has `.name`, `.description`, `.arguments`
4. Iterates `command_group_table` — each entry has group help text
5. Handles: callable descriptions, callable choices, non-string help values, broken extensions
6. Outputs JSON to stdout

### Docker Scraping (`Dockerfile` + `compose.yaml`)

```dockerfile
ARG VERSION=2.71.0
FROM mcr.microsoft.com/azure-cli:${VERSION}
COPY scripts/dump_command_table.py /
CMD ["python", "/dump_command_table.py"]
```

Dependabot tracks the `docker` ecosystem and opens PRs when new az CLI versions are published. The dependabot workflow then runs the full pipeline: Docker scrape → spec generation → code generation → commit → auto-merge.

## Gotchas

- **`azcli_generated.go` is generated**: Always run `go generate ./cmd/carapace-az` after adding, removing, or renaming YAML spec files. Never hand-edit it.
- **YAML spec files are embedded**: The `//go:embed *.yaml` directive in `azcli.go` embeds all `.yaml` files in the `azcli/` directory at build time. Adding/removing YAML files requires rebuilding.
- **Specs are loaded lazily**: Group specs are only unmarshaled in `PreRun` (when the user invokes that group), not at startup. This keeps startup fast given the large number of command groups.
- **Bridge arg passthrough**: `bridge.go` uses `os.Args[4:]` to reconstruct the completion context — this is fragile and depends on the caller's argument structure.
- **Bridged completions**: Most actual completion logic lives in `carapace-bridge` (external). This repo defines the command structure and delegates to it.
- **pflag replacement**: The `replace github.com/spf13/pflag => github.com/carapace-sh/carapace-pflag v1.1.0` directive in `go.mod` ensures carapace's patched pflag is used.
- **No CamelCase conversion needed**: az CLI already uses kebab-case flag names (e.g. `--resource-group`), unlike AWS which needs CamelCaseToDash conversion.
- **Nested command hierarchy**: az CLI has deeply nested groups (e.g., `az network vnet subnet create`). Specs use nested `commands` in YAML to represent this — one file per top-level group.
- **Extensions not scraped**: The Docker scrape only captures core command modules. az CLI extensions (e.g., `az devops`, `az kusto`, `az ml`) are not included. Install extensions in the Dockerfile before scraping to include them.
- **Docker/compose**: The `Dockerfile` and `compose.yaml` are for generating the command table JSON, not for running the completer.
- **Man page docs are separate from completion specs**: `man/cmd/az/` contains documentation (from knack `long-summary`) and is generated with `carapace-spec-az --stdout` + `carapace-man update`. Completion specs in `cmd/carapace-az/cmd/azcli/` are generated with `--no-doc` to keep them small.
- **Man page enrichment from learn.microsoft.com**: The automated pipeline generates baseline docs from knack help text. On top of that, `documentation.flag` entries and `examples:` in ~343 files were manually enriched by scraping the official Azure CLI reference pages at `https://learn.microsoft.com/en-us/cli/azure/<group>`. This is a one-time manual step — `carapace-man update` preserves existing documentation on re-generation, so the enriched docs survive automated pipeline runs. To re-enrich after major az CLI updates, fetch the relevant group pages, extract command descriptions, flag docs, and examples, and write them into the `man/cmd/az/*.yaml` files.
- **Partial doc coverage**: knack `long-summary` is not present for all commands. Groups with sub-pages (e.g. `network`, `storage`, `keyvault`) only have top-level commands documented on their main page; sub-commands live on separate pages (e.g. `/cli/azure/network/vnet`). A full enrichment would require fetching those sub-pages as well.
- **Go version**: `go.mod` specifies `go 1.25.0`.

## Maintenance Workflow (Dependabot Updates)

When a new az CLI Docker image version is published, dependabot opens a PR. The dependabot workflow automatically:

1. Builds the Docker image and runs `dump_command_table.py` to produce `az_commands.json`
2. Runs `carapace-spec-az --no-doc` to convert JSON to completion YAML specs (no documentation)
3. Installs `carapace-man` and runs `carapace-spec-az --stdout` + `carapace-man update` to generate man page specs in `man/cmd/az/`
4. Runs `go generate` to update `azcli_generated.go`
5. Commits and pushes changes, then auto-merges the PR

### Manual Update

To manually update specs for a specific az CLI version:

```sh
# Update VERSION in Dockerfile and compose.yaml
# Then run the full pipeline:
docker compose build && docker compose run --rm az > az_commands.json
go run -C cmd/carapace-spec-az . --target cmd/carapace-az/cmd/azcli --no-doc az_commands.json
go run -C cmd/carapace-spec-az . --stdout az_commands.json > /tmp/az-full-spec.yaml
carapace-man update /tmp/az-full-spec.yaml man/cmd/az
go generate ./cmd/carapace-az
go build ./cmd/carapace-az  # verify build
```
