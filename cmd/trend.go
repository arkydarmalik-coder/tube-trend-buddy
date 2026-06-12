package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/llm"
	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/prompts"
)

var trendCmd = &cobra.Command{
	Use:   "trend",
	Short: "Detect rising niches + trending topics",
	Long:  "Uses LLM knowledge + (optional) YouTube Data API context to surface rising niches.",
	Example: `  ttb trend --region ID --period 7d
  ttb trend --region US --category tech --count 15`,
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
	system, user := prompts.Trend(flagRegion, flagCategory, flagPeriod, cfg.Lang, cfg.Count)
	out, err := client.Complete(ctx, system, user)
	if err != nil {
		return err
	}
	if cfg.JSON {
		fmt.Println(out)
		return nil
	}
	fmt.Printf("Rising niches in %q (category=%s, period=%s):\n\n", flagRegion, flagCategory, flagPeriod)
	for _, line := range splitList(out) {
		fmt.Println("  > " + line)
	}
	return nil
}
