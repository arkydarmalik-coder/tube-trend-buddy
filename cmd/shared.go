// Package shared holds CLI helpers shared across all 8 subcommands:
// config resolution, output splitting, and the LLM client lifecycle.
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/llm"
)

// shared is a tiny holder so subcommands can chain helpers in a stable order.
var shared = &sharedDeps{}

type sharedDeps struct{}

// applyRootFlags is a no-op kept for forward compatibility (per-command flag
// overrides will hook in here).
func (s *sharedDeps) applyRootFlags(_ *cobra.Command) {}

// resolveConfig merges: cobra flags > TTB_* env > provider-specific env > defaults.
func resolveConfig(cmd *cobra.Command) (llm.Config, error) {
	cfg := llm.Config{
		Provider: pickString(flagProvider, "TTB_PROVIDER", "openrouter"),
		Model:    pickString(flagModel, "TTB_MODEL", ""),
		Lang:     pickString(flagLang, "TTB_LANG", "en"),
		Count:    flagCount,
		JSON:     flagJSON,
		Timeout:  90 * time.Second,
	}

	provider := strings.ToLower(cfg.Provider)
	cfg.APIKey, cfg.BaseURL, cfg.Model = pickProviderDefaults(provider, cfg.Model)

	if cfg.Model == "" {
		return cfg, fmt.Errorf("no model resolved for provider %q (set --model or TTB_MODEL)", provider)
	}
	if cfg.APIKey == "" {
		return cfg, fmt.Errorf("no API key for provider %q (set TTB_API_KEY or provider-specific env)", provider)
	}
	return cfg, nil
}

func pickString(flagVal, envName, def string) string {
	if v := strings.TrimSpace(flagVal); v != "" {
		return v
	}
	if v := strings.TrimSpace(os.Getenv(envName)); v != "" {
		return v
	}
	return def
}

func pickProviderDefaults(provider, model string) (key, baseURL, modelOut string) {
	switch provider {
	case "openrouter":
		key = pickString("", "TTB_API_KEY", os.Getenv("OPENROUTER_API_KEY"))
		baseURL = "https://openrouter.ai/api/v1"
		if model == "" {
			model = "anthropic/claude-sonnet-4.6"
		}
	case "naraya":
		key = pickString("", "TTB_API_KEY", os.Getenv("NARAYA_API_KEY"))
		baseURL = "https://router.naraya.ai/v1"
		if model == "" {
			model = "MiniMax-M3"
		}
	case "huggingface", "hf":
		key = pickString("", "TTB_API_KEY", os.Getenv("HUGGINGFACE_HUB_TOKEN"))
		baseURL = "https://router.huggingface.co/v1"
		if model == "" {
			model = "Qwen/Qwen3.5-35B-A3B"
		}
	case "gemini":
		key = pickString("", "TTB_API_KEY", os.Getenv("GOOGLE_API_KEY"))
		// Gemini has an OpenAI-compat shim on some gateways; if user wants native
		// Gemini API, override --base-url via env. Default to the OpenAI-compat
		// route so the same client code path works.
		baseURL = pickString("", "TTB_BASE_URL", "https://generativelanguage.googleapis.com/v1beta/openai/")
		if model == "" {
			model = "gemini-2.5-pro"
		}
	case "ollama":
		key = "ollama" // Ollama doesn't need a real key
		baseURL = pickString("", "TTB_BASE_URL", "http://localhost:11434/v1")
		if model == "" {
			model = "llama3.2"
		}
	default:
		baseURL = pickString("", "TTB_BASE_URL", "")
		key = pickString("", "TTB_API_KEY", "")
	}
	return key, baseURL, model
}

// splitList takes raw LLM output (which is usually a list with bullets or
// numbered prefixes) and returns one cleaned item per line.
func splitList(s string) []string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		// strip common prefixes
		for _, prefix := range []string{"- ", "* ", "• ", "· "} {
			if strings.HasPrefix(ln, prefix) {
				ln = strings.TrimPrefix(ln, prefix)
				break
			}
		}
		// strip "1. ", "12) " etc
		if i := indexNumbered(ln); i >= 0 {
			ln = ln[i:]
		}
		ln = strings.TrimSpace(ln)
		if ln != "" {
			out = append(out, ln)
		}
	}
	return out
}

func indexNumbered(s string) int {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			continue
		}
		if c == '.' || c == ')' || c == ':' {
			if i+1 < len(s) && s[i+1] == ' ' {
				return i + 2
			}
		}
		break
	}
	return -1
}
