// Package cmd wires up the cobra CLI tree. Each file in this package adds
// one subcommand that maps 1:1 to one of the 8 main creator features.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagProvider string
	flagModel    string
	flagLang     string
	flagJSON     bool
	flagCount    int
	flagYolo     bool
)

var rootCmd = &cobra.Command{
	Use:   "ttb",
	Short: "Tube Trend Buddy - 8-in-1 YouTube creator toolkit",
	Long: `Tube Trend Buddy (ttb) helps YouTube creators ship better videos faster.

8 features, one binary, multi-LLM provider:
  title      AI title generator (CTR-optimized)
  tags       SEO tag optimizer (rank + traffic)
  trend      Trend detector (rising niches)
  niche      Niche + channel analyzer
  describe   Description + hashtag writer
  thumbnail  Thumbnail concept generator
  monetize   Monetization strategy advisor
  calendar   Content calendar planner

Provider priority: --provider flag > TTB_PROVIDER env > openrouter.
API key: --api-key flag > TTB_API_KEY env > provider-specific env.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagProvider, "provider", "",
		"LLM provider: openrouter | naraya | gemini | huggingface | ollama")
	rootCmd.PersistentFlags().StringVar(&flagModel, "model", "",
		"LLM model id (provider-specific, e.g. anthropic/claude-sonnet-4.6)")
	rootCmd.PersistentFlags().StringVar(&flagLang, "lang", "en",
		"Output language: en | id | es | etc.")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false,
		"Print raw LLM response as JSON")
	rootCmd.PersistentFlags().IntVar(&flagCount, "count", 10,
		"Number of suggestions to generate")
	rootCmd.PersistentFlags().BoolVar(&flagYolo, "yolo", false,
		"Skip safety disclaimer on the first call")
}
