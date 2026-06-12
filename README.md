# Tube Trend Buddy (`ttb`)

**8-in-1 YouTube creator toolkit, one Go binary, multi-LLM.**

`ttb` is a small command-line app that gives YouTube creators a fast, repeatable
way to ship better videos: titles, SEO tags, trend detection, niche audits,
descriptions, thumbnail concepts, monetization plans, and content calendars.
No UI, no database, no telemetry. One executable, your API key, done.

## 8 features (one subcommand each)

| Subcommand | What it does |
|---|---|
| `ttb title` | Generate CTR-optimized titles for a niche (curiosity, clarity, no clickbait) |
| `ttb tags` | SEO tag mix (broad + mid + long tail) tuned for discoverability |
| `ttb trend` | Detect rising niches + suggested content angles |
| `ttb niche` | Audit a YouTube channel: positioning, gaps, adjacent niches, 90-day experiments |
| `ttb describe` | Write description + hashtags + (optional) CTA + (optional) chapters |
| `ttb thumbnail` | Thumbnail concept with text overlay, mood, palette, copy-pasteable image-gen prompt |
| `ttb monetize` | Tiered monetization plan (ads, sponsorships, products, affiliate, memberships) |
| `ttb calendar` | Full-month content calendar with mixed formats |

## Install

```bash
go install github.com/arkydarmalik-coder/tube-trend-buddy@latest
```

Or build from source:

```bash
git clone https://github.com/arkydarmalik-coder/tube-trend-buddy
cd tube-trend-buddy
make build
# cross-compile a Windows .exe
make build-windows
# cross-compile for mac/linux
make build-all
```

## Quick start

```bash
# 1. Set your LLM API key (any OpenAI-compatible provider)
export TTB_API_KEY="sk-or-v1-..."          # OpenRouter example
export TTB_PROVIDER="openrouter"           # or: naraya | gemini | huggingface | ollama
export TTB_MODEL="anthropic/claude-sonnet-4.6"  # optional, provider has a default

# 2. Run a feature
ttb title --niche "AI tools" --count 10 --lang en
ttb tags  --title "How to Learn AI in 30 Days" --count 30
ttb trend --region ID --category tech --period 7d
ttb niche --channel @mkbhd --deep
ttb describe --title "I Built an AI Agent" --cta
ttb thumbnail --title "I Built an AI Agent" --mood shocked --face 1 --colors "yellow,black,red"
ttb monetize --niche "AI tutorial" --subs 5000 --region ID
ttb calendar --month 7 --niche "AI tools" --frequency 2/week

# 3. Pipe-friendly output
ttb tags --title "..." --json | jq '.[]'
ttb title --niche "AI" --count 5 --json
```

## Provider matrix

`ttb` speaks OpenAI-compatible `/v1/chat/completions`. Set `TTB_PROVIDER` and
`TTB_API_KEY` (or rely on the provider-specific env var below). You can override
`TTB_BASE_URL` to point at any compatible gateway.

| Provider | `TTB_PROVIDER` | Default base URL | Default model | Auto-picked env key |
|---|---|---|---|---|
| OpenRouter | `openrouter` | `https://openrouter.ai/api/v1` | `anthropic/claude-sonnet-4.6` | `OPENROUTER_API_KEY` |
| Naraya | `naraya` | `https://router.naraya.ai/v1` | `MiniMax-M3` | `NARAYA_API_KEY` |
| HuggingFace Router | `huggingface` | `https://router.huggingface.co/v1` | `Qwen/Qwen3.5-35B-A3B` | `HUGGINGFACE_HUB_TOKEN` |
| Google Gemini (compat) | `gemini` | `https://generativelanguage.googleapis.com/v1beta/openai/` | `gemini-2.5-pro` | `GOOGLE_API_KEY` |
| Ollama (local) | `ollama` | `http://localhost:11434/v1` | `llama3.2` | (none, uses `ollama`) |

> **Tip:** `TTB_API_KEY` always wins over the provider-specific env. Useful for
> running multiple `ttb` profiles with different keys.

## Project layout

```
tube-trend-buddy/
  main.go                    # entrypoint
  cmd/                       # 1 file per subcommand (8 features + shared helpers)
    root.go                  # cobra root, global flags
    shared.go                # config resolution, output splitter
    title.go                 # 1. title generator
    tags.go                  # 2. SEO tag optimizer
    trend.go                 # 3. trend detector
    niche.go                 # 4. niche analyzer
    describe.go              # 5. description + hashtag writer
    thumbnail.go             # 6. thumbnail concept
    monetize.go              # 7. monetization advisor
    calendar.go              # 8. content calendar
  internal/
    llm/client.go            # OpenAI-compatible HTTP client (no SDK)
    prompts/prompts.go       # 8 prompt templates (pure functions)
  configs/
    default.yaml             # sample config (env vars override)
  Makefile                   # build, run, cross-compile
  README.md
  LICENSE                    # MIT
```

## Design notes

- **No SDKs.** The LLM client is ~120 lines of `net/http`. Easier to audit and
  trivially swappable.
- **Prompts are pure functions.** `internal/prompts/prompts.go` exposes 8
  `(system, user) -> string` pairs. Edit them in one place.
- **No state.** No SQLite, no Redis, no telemetry. `--json` is the only output
  switch you need.
- **Cross-compile friendly.** `make build-all` produces `ttb-linux-amd64`,
  `ttb-darwin-arm64`, and `ttb-windows-amd64.exe`.

## Known limitations

- No streaming (single-shot completions only). Easy to add later.
- No YouTube Data API integration yet (the `trend` and `niche` commands rely on
  the LLM's own knowledge). Pull requests welcome.
- No persistent history. Re-run when you want a fresh take.

## License

MIT. See `LICENSE`.
