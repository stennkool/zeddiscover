# AGENTS.md

## What this is

Go CLI that syncs Zed editor's `~/.config/zed/settings.json` with live provider `/models` endpoints. Reads `language_models.openai_compatible` entries, fetches model lists, filters to text-only models, and writes back `available_models`.

## Commands

- `go run .` — sync all providers (writes to disk)
- `go run . --dry-run` — preview changes only
- `go run . --provider <name>` — sync a single provider
- `go run . --config <path>` — use a custom config path
- `task` / `task default` — dry-run sync (Taskfile shorthand)
- `task build` — build `zeddiscover` binary
- `task install` — install binary to `~/.local/bin`

## Architecture

```
main.go            CLI flags, wires dependencies
internal/config    Read/write Zed settings.json (handles JSONC trailing commas, creates .bak backup)
internal/model     Domain types (APIModel → AvailableModel mapping, text-only filter)
internal/provider  Fetcher interface + OpenRouter-style implementation
internal/sync      Runner: iterates providers, fetches models, updates config
```

- `internal/` is intentional: nothing is importable externally.
- Only one fetcher implementation exists (`OpenRouterFetcher`); adding a new provider requires a new `provider.Fetcher` implementation.
- Config backup: `Write()` copies the existing file to `<path>.bak` before overwriting.

## Testing

No tests exist yet. When adding tests, `go test ./...` from repo root.

## Gotchas

- Config parsing strips JSONC trailing commas via regex — does not handle `//` line comments.
- Default config path requires `~/.config/zed/settings.json` to exist at runtime.
