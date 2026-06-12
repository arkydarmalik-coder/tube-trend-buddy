package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/llm"
	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/prompts"
)

var monetizeCmd = &cobra.Command{
	Use:   "monetize",
	Short: "Monetization strategy advisor for a channel",
	Long:  "Recommends a tiered plan: ads, sponsorships, products, affiliate, memberships.",
	Example: `  ttb monetize --niche "AI tutorial" --subs 5000 --region ID
  ttb monetize --niche "..." --subs 50000 --focus sponsorships`,
	RunE: runMonetize,
}

var (
	flagSubs  int
	flagFocus string
)

func init() {
	rootCmd.AddCommand(monetizeCmd)
	monetizeCmd.Flags().StringVar(&flagNiche, "niche", "", "Channel niche (required)")
	monetizeCmd.Flags().IntVar(&flagSubs, "subs", 1000, "Approx subscriber count")
	monetizeCmd.Flags().StringVar(&flagRegion, "region", "US", "Primary audience region")
	monetizeCmd.Flags().StringVar(&flagFocus, "focus", "all", "Focus: ads | sponsorships | products | affiliate | memberships | all")
	_ = monetizeCmd.MarkFlagRequired("niche")
}

func runMonetize(cmd *cobra.Command, args []string) error {
	cfg, err := resolveConfig(cmd)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client, err := llm.New(cfg)
	if err != nil {
		return err
	}
	system, user := prompts.Monetize(flagNiche, flagSubs, flagRegion, flagFocus, cfg.Lang)
	out, err := client.Complete(ctx, system, user)
	if err != nil {
		return err
	}
	if cfg.JSON {
		fmt.Println(out)
		return nil
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Monetization plan for %q (~%d subs, %s, focus=%s):\n\n",
		flagNiche, flagSubs, flagRegion, flagFocus)
	for _, line := range splitList(out) {
		fmt.Fprintln(tw, "  - "+line)
	}
	return tw.Flush()
}
