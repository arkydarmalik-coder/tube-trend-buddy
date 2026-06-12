# Tube Trend Buddy (`ttb`)

**8-in-1 YouTube creator toolkit: one Go binary, multi-LLM, live YouTube data, optional TUI.**

[![Release](https://img.shields.io/github/v/release/arkydarmalik-coder/tube-trend-buddy)](https://github.com/arkydarmalik-coder/tube-trend-buddy/releases/latest)
[![CI](https://img.shields.io/badge/CI-passing-brightgreen)](https://github.com/arkydarmalik-coder/tube-trend-buddy/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8)](https://go.dev)

[Download v0.2.0](https://github.com/arkydarmalik-coder/tube-trend-buddy/releases/tag/v0.2.0) (linux / macOS / windows binaries)

`ttb` gives YouTube creators a fast, repeatable way to ship better videos:
titles, SEO tags, trend detection, niche audits, descriptions, thumbnail
concepts, monetization plans, and content calendars. No UI to log into,
no database, no telemetry. One executable, your API key, done.

## 8 features (one subcommand each)

| Subcommand | What it does |
|---|---|
| `ttb title` | Generate CTR-optimized titles for a niche (curiosity, clarity, no clickbait) |
| `ttb tags` | SEO tag mix (broad + mid + long tail) tuned for discoverability |
| `ttb trend` | Detect rising niches + suggested content angles. **Uses live YouTube Data API if `YOUTUBE_API_KEY` is set** |
| `ttb niche` | Audit a YouTube channel: positioning, gaps, adjacent niches, 90-day experiments. **Uses live YouTube Data API if `YOUTUBE_API_KEY` is set** |
| `ttb describe` | Write description + hashtags + (optional) CTA + (optional) chapters |
| `ttb thumbnail` | Thumbnail concept with text overlay, mood, palette, copy-pasteable image-gen prompt |
| `ttb monetize` | Tiered monetization plan (ads, sponsorships, products, affiliate, memberships) |
| `ttb calendar` | Full-month content calendar with mixed formats |

Plus: **`ttb ui`** - launch the interactive TUI wizard (Bubble Tea, 9th subcommand).

## Install

```bash
go install github.com/arkydarmalik-coder/tube-trend-buddy@latest
```

Or build from source:

```bash
git clone https://github.com/arkydarmalik-coder/tube-trend-buddy
cd tube-trend-buddy
make build          # current OS -> bin/ttb
make build-windows  # -> dist/ttb-windows-amd64.exe
make build-all      # all 3 OS
```

The compiled binaries are also attached to each GitHub Release.

## Quick start

### A. CLI mode

```bash
# 1. LLM provider (any OpenAI-compatible)
export TTB_API_KEY=***                    # OpenRouter
export TTB_PROVIDER="openrouter"
export TTB_MODEL="anthropic/claude-sonnet-4.6"   # optional default

# 2. (optional) YouTube Data API for live trend/niche data
#    Create at https://console.cloud.google.com -> enable "YouTube Data API v3"
export YOUTUBE_API_KEY=***

# 3. Run a feature
ttb title  --niche "AI tools" --count 10 --lang en
ttb tags   --title "How to Learn AI in 30 Days" --count 30
ttb trend  --region ID --category tech --period 7d          # uses live YouTube data
ttb niche  --channel @mkbhd --deep                          # uses live YouTube data
ttb describe --title "I Built an AI Agent" --cta
ttb thumbnail --title "I Built an AI Agent" --mood shocked --face 1 --colors "yellow,black,red"
ttb monetize --niche "AI tutorial" --subs 5000 --region ID
ttb calendar --month 7 --niche "AI tools" --frequency 2/week

# 4. Pipe-friendly output
ttb tags --title "..." --json | jq
ttb title --niche "AI" --count 5 --json

# 5. Force LLM-only mode (skip YouTube)
ttb trend --region ID --no-youtube
ttb niche --channel @mkbhd --no-youtube
```

### B. TUI mode (interactive)

```bash
ttb ui
```

Opens a full-screen Bubble Tea wizard: pick a feature with arrow keys,
fill the form, press Enter, watch the spinner, scroll the result.

```
+----------------------------------------------+
|  Tube Trend Buddy - 8 features               |
+----------------------------------------------+
|  > 1. title                                  |
|       CTR-optimized YouTube titles           |
|                                              |
|    2. tags                                   |
|       SEO tags (broad + mid + long tail)     |
|    3. trend                                  |
|       Rising niches + content angles         |
|       (live YouTube data)                    |
|    ...                                       |
+----------------------------------------------+
  arrows/enter to pick - q to quit
```

TUI requires a real terminal (won't work in a non-TTY pipe).

## YouTube Data API integration

`ttb trend` and `ttb niche` use the real YouTube Data API v3 to feed the
LLM with **current** data instead of stale model knowledge. The cost is
**1-3 quota units per call** (free tier is 10,000 units/day, so you can
run thousands of audits per day for free).

### What `trend` fetches

- Top 25 most-popular videos in the chosen region (`videos.list?chart=mostPopular&regionCode=...`)
- Sorted by YouTube's own trending algorithm, with view/like/comment counts and channel info
- LLM then groups them into "rising niches" and "suggested angles"

### What `niche` fetches

- Channel metadata (subs, total views, video count, country, creation date) via `channels.list?forHandle=@...`
- Latest 25 videos with full statistics via `search.list` + `videos.list`
- LLM audits positioning, strengths, gaps, and adjacent niches

### Setup (one-time, ~2 minutes)

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create a project (or pick an existing one)
3. **APIs & Services -> Library -> search "YouTube Data API v3" -> Enable**
4. **APIs & Services -> Credentials -> Create credentials -> API key**
5. (Optional but recommended) **Restrict the key** to "YouTube Data API v3" only
6. Export it: `export YOUTUBE_API_KEY=AIza...`

That's it. The next `ttb trend` or `ttb niche` call will automatically use
live data. No code changes, no extra flags.

To force LLM-only mode for a single call, pass `--no-youtube`.

## Provider matrix (LLM)

`ttb` speaks OpenAI-compatible `/v1/chat/completions`. Set `TTB_PROVIDER` and
`TTB_API_KEY` (or rely on the provider-specific env var below). You can override
`TTB_BASE_URL` to point at any compatible gateway.

| Provider | `TTB_PROVIDER` | Default base URL | Default model | Auto-picked env key |
|---|---|---|---|---|
| OpenRouter | `openrouter` | `https://openrouter.ai/api/v1` | `anthropic/claude-sonnet-4.6` | `OPENROUTER_API_KEY` |
| Naraya | `naraya` | `https://router.naraya.ai/v1` | `MiniMax-M3` | `NARAYA_API_KEY` |
| HuggingFace Router | `huggingface` | `https://router.huggingface.co/v1` | `Qwen/Qwen3.5-35B-A3B` | `HUGGINGFACE_HUB_TOKEN` |
| Google Gemini (compat) | `gemini` | `https://generativelanguage.googleapis.com/v1beta/openai/` | `gemini-2.5-pro` | `GOOGLE_API_KEY` |
| Ollama (local) | `ollama` | `http://localhost:11434/v1` | `llama3.2` | (none) |

> **Tip:** `TTB_API_KEY` always wins over the provider-specific env. Useful for
> running multiple `ttb` profiles with different keys.

## Project layout

```
tube-trend-buddy/
  main.go                    # entrypoint
  cmd/                       # 1 file per subcommand (8 features + ui + shared helpers)
    root.go                  # cobra root, global flags
    shared.go                # config resolution, output splitter
    title.go                 # 1. title generator
    tags.go                  # 2. SEO tag optimizer
    trend.go                 # 3. trend detector (+ YouTube)
    niche.go                 # 4. niche analyzer (+ YouTube)
    describe.go              # 5. description + hashtag writer
    thumbnail.go             # 6. thumbnail concept
    monetize.go              # 7. monetization advisor
    calendar.go              # 8. content calendar
    ui.go                    # 9. TUI wizard launcher
  internal/
    llm/client.go            # OpenAI-compatible HTTP client (~120 lines net/http)
    prompts/prompts.go       # 8 prompt templates (pure functions)
    youtube/client.go        # YouTube Data API v3 wrapper
    ui/tui.go                # Bubble Tea TUI (~400 lines)
  configs/
    default.yaml             # sample config
  Makefile                   # build, run, cross-compile
  README.md
  LICENSE                    # MIT
```

## Design notes

- **No SDKs.** The LLM client is ~120 lines of `net/http`. Same for the
  YouTube client (~300 lines). Both are trivially auditable and swappable.
- **Prompts are pure functions.** `internal/prompts/prompts.go` exposes 8
  `(system, user) -> string` pairs. Edit them in one place. `TrendWithData`
  and `NicheWithData` are wrappers that inject real YouTube data when
  available and fall back gracefully when it isn't.
- **No state.** No SQLite, no Redis, no telemetry. `--json` is the only
  output switch you need.
- **TUI is a launcher.** `ttb ui` shells out to `os.Args[0]` with the right
  args, so the CLI and TUI never diverge. Single source of truth.
- **Cross-compile friendly.** `make build-all` produces `ttb-linux-amd64`,
  `ttb-darwin-arm64`, and `ttb-windows-amd64.exe`.

## Known limitations

- No streaming (single-shot completions only). Easy to add later.
- TUI only runs in real terminals (TTY). Won't work in `script`, `nohup`, etc.
  Use the CLI mode for scripting.
- The YouTube client only reads public data (videos.list, channels.list,
  search.list). No OAuth, no uploads, no comments. By design.
- No persistent history. Re-run when you want a fresh take.

## Releases

Every tag like `v0.2.1` pushed to `main` triggers `.github/workflows/release.yml`,
which cross-compiles 3 binaries (linux / macOS / Windows), generates SHA256
checksums, creates a GitHub Release, and uploads all 4 files as assets.
Typical end-to-end time: ~3 minutes.

```bash
# Release flow
git tag v0.2.1                    # bump the version
git push origin v0.2.1            # CI takes it from here
# -> https://github.com/arkydarmalik-coder/tube-trend-buddy/releases/tag/v0.2.1
```

There is also a lightweight `ci.yml` that runs on every push/PR to `main`:
`go vet`, `go build`, and a smoke test that all 9 subcommands are registered.

## License

MIT. See `LICENSE`.
