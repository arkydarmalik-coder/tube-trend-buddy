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

var nicheCmd = &cobra.Command{
	Use:   "niche",
	Short: "Audit a YouTube channel + find niche opportunities",
	Long:  "Reviews positioning, content gaps, and adjacent niches worth expanding into.",
	Example: `  ttb niche --channel @mkbhd
  ttb niche --channel @mrbeast --deep`,
	RunE: runNiche,
}

var (
	flagChannel string
	flagDeep    bool
)

func init() {
	rootCmd.AddCommand(nicheCmd)
	nicheCmd.Flags().StringVar(&flagChannel, "channel", "", "Channel handle or URL (required)")
	nicheCmd.Flags().BoolVar(&flagDeep, "deep", false, "Deeper analysis (slower, more tokens)")
	_ = nicheCmd.MarkFlagRequired("channel")
}

func runNiche(cmd *cobra.Command, args []string) error {
	cfg, err := resolveConfig(cmd)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client, err := llm.New(cfg)
	if err != nil {
		return err
	}
	system, user := prompts.Niche(flagChannel, flagDeep, cfg.Lang)
	out, err := client.Complete(ctx, system, user)
	if err != nil {
		return err
	}
	if cfg.JSON {
		fmt.Println(out)
		return nil
	}
	fmt.Printf("Niche audit for %q:\n\n", flagChannel)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, line := range splitList(out) {
		fmt.Fprintln(tw, "  - "+line)
	}
	return tw.Flush()
}
