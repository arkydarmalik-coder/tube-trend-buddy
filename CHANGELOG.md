# Changelog

## 0.2.0 (2026-06-12)

Added 2 big features on top of v0.1.0:

- **TUI wizard (`ttb ui`)** - full-screen interactive picker built on Bubble
  Tea + Lip Gloss + Bubbles. 8 menu items, per-feature form, spinner, scrollable
  result viewport. Falls back to the existing CLI for actual execution
  (shells out to `os.Args[0]`, so there is zero logic duplication).
- **YouTube Data API v3 integration** for `trend` and `niche`. When
  `YOUTUBE_API_KEY` is set, those commands now fetch the **current** top 25
  trending videos (by region) or the latest 25 videos from a specific channel,
  with full statistics, and inject them into the LLM prompt as context.
  Add `--no-youtube` to force the old LLM-only behavior.

Internal changes:

- New package `internal/youtube` (HTTP wrapper around YouTube Data API v3,
  ~300 lines, no SDK).
- New package `internal/ui` (Bubble Tea TUI, ~400 lines).
- New `TrendWithData` and `NicheWithData` prompt variants in
  `internal/prompts/prompts.go`.
- New `pickYouTubeClient()` helper in `cmd/trend.go` (shared with `niche.go`).
- Binary size: 6.9 MB -> 9-10 MB (added Bubble Tea + YouTube client deps).
- New env var: `YOUTUBE_API_KEY` (and `TTB_YOUTUBE_API_KEY` override).

## 0.1.0 (2026-06-12)

Initial release. 8 features in one Go binary:

- `ttb title` - CTR-optimized title generator
- `ttb tags` - SEO tag optimizer
- `ttb trend` - Trend detector
- `ttb niche` - Channel/niche analyzer
- `ttb describe` - Description + hashtag writer
- `ttb thumbnail` - Thumbnail concept generator
- `ttb monetize` - Monetization strategy advisor
- `ttb calendar` - Content calendar planner

Multi-provider LLM support: OpenRouter, Naraya, HF Router, Gemini (compat), Ollama.
