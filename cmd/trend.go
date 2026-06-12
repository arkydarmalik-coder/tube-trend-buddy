package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/llm"
	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/prompts"
	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/youtube"
)

var trendCmd = &cobra.Command{
	Use:   "trend",
	Short: "Detect rising niches + trending topics",
	Long:  "Uses LLM knowledge + (optional) YouTube Data API v3 for live trending context.",
	Example: `  ttb trend --region ID --period 7d
  ttb trend --region US --category tech --count 15
  ttb trend --region ID --no-youtube              # LLM-only mode`,
	RunE: runTrend,
}

var (
	flagRegion   string
	flagCategory string
	flagPeriod   string
)

func init() {
	rootCmd.AddCommand(trendCmd)
	trendCmd.Flags().StringVar(&flagRegion, "region", "US", "ISO country code or region")
	trendCmd.Flags().StringVar(&flagCategory, "category", "general", "Category hint: tech | gaming | finance | etc.")
	trendCmd.Flags().StringVar(&flagPeriod, "period", "7d", "Lookback period: 1d | 7d | 30d | 90d")
	trendCmd.Flags().BoolVar(&flagNoYouTube, "no-youtube", false, "Force LLM-only mode (skip YouTube Data API)")
}

func runTrend(cmd *cobra.Command, args []string) error {
	cfg, err := resolveConfig(cmd)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client, err := llm.New(cfg)
	if err != nil {
		return err
	}
	ytClient := pickYouTubeClient()
	if flagNoYouTube {
		ytClient = nil
	}
	system, user := prompts.TrendWithData(ctx, ytClient, flagRegion, flagCategory, flagPeriod, cfg.Lang, cfg.Count)
	out, err := client.Complete(ctx, system, user)
	if err != nil {
		return err
	}
	if cfg.JSON {
		fmt.Println(out)
		return nil
	}
	header := fmt.Sprintf("Rising niches in %q (category=%s, period=%s", flagRegion, flagCategory, flagPeriod)
	if ytClient != nil {
		header += ", youtube=live"
	} else {
		header += ", youtube=off"
	}
	header += "):"
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)))
	for _, line := range splitList(out) {
		fmt.Println("  > " + line)
	}
	return nil
}

// pickYouTubeClient returns a *youtube.Client if YOUTUBE_API_KEY (or TTB_YOUTUBE_API_KEY)
// is set, otherwise nil. Callers should pass nil to prompt builders to signal LLM-only mode.
func pickYouTubeClient() *youtube.Client {
	key := pickString("", "TTB_YOUTUBE_API_KEY", "")
	if key == "" {
		key = os.Getenv("YOUTUBE_API_KEY")
	}
	if key == "" {
		return nil
	}
	return youtube.New(key)
}
