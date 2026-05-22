# zeddiscover

**Automatically sync the Zed editor's language model list from live provider API endpoints.**

---

## The problem

The [Zed editor](https://zed.dev) lets you connect to any OpenAI-compatible API under
`language_models.openai_compatible`. Great — but every model you want to use must be
declared by hand in `~/.config/zed/settings.json` as an `available_models` entry with
its name, context length, capabilities, and other metadata.

Kilo, OpenRouter, Synthetic, and other gateways serve **hundreds** of models that rotate
daily.  Keeping that list up to date by hand is:

- **Tedious** — copy-pasting model IDs one by one.
- **Error-prone** — wrong token limits break the model selector UI.
- **Stale** — new models launch; deprecated ones disappear.

`zeddiscover` automates this away.  It reads your existing providers, hits each one's
`/models` endpoint, filters to text-only LLMs, and writes a clean `available_models`
array.  Run it on a cron, and your model picker stays current forever.

## Install

```sh
git clone https://github.com/stenn/zeddiscover.git
cd zeddiscover

# Build & install to ~/.local/bin
task install

# Or just build the binary
task build
```

Requires **Go 1.21+**.

## Usage

```sh
# Sync all OpenAI-compatible providers
zeddiscover

# Preview what would change (does not touch disk)
zeddiscover --dry-run

# Sync a single provider by name
zeddiscover --provider kilo

# Point at a custom config file
zeddiscover --config /path/to/settings.json
```

You can also use `go run .` from the repository:

```sh
go run .                      # sync all
go run . --dry-run            # preview
go run . --provider kilo      # single provider
task                          # shortcut: dry-run all
```

### Taskfile shortcuts

| Command         | What it does                          |
| --------------- | ------------------------------------- |
| `task`          | Dry-run sync (`go run . --dry-run`)   |
| `task build`    | Build the `zeddiscover` binary        |
| `task install`  | Build + install to `~/.local/bin`     |
| `task run --`   | Run with arguments (e.g. `--provider`)|

## How it works

1. **Read** — parses `settings.json`, finds every `language_models.openai_compatible`
   entry that has an `api_url`.
2. **Fetch** — calls `<api_url>/models` to retrieve the live model list.
3. **Filter** — keeps only models that produce **text output** (discards image, audio,
   and video generators — those aren't useful for chatting).
4. **Back up** — copies `settings.json` → `settings.json.bak` before touching the file.
5. **Write** — replaces each provider's `available_models` array and saves.

## Configuration in Zed

Define one or more OpenAI-compatible providers in your Zed `settings.json`:

```json
"language_models": {
  "openai_compatible": {
    "kilo": {
      "api_url": "https://api.kilo.ai/api/gateway"
    },
    "synthetic": {
      "api_url": "https://api.synthetic.new/openai/v1"
    }
  }
}
```

The `available_models` array will be populated automatically the next time you run
`zeddiscover`.  No manual model entries needed.

## Project layout

```
main.go              CLI flags, wires dependencies
internal/
  config/            Read, write, and backup settings.json (handles JSONC)
  model/             Domain types, text-only filter, API → Zed model mapping
  provider/          Fetcher interface + OpenRouter/Kilo/Synthetic implementation
  sync/              Runner: iterates providers, fetches models, writes config
```

## License

MIT
