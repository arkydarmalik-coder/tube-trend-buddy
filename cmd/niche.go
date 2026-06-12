package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/llm"
	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/prompts"
)

var nicheCmd = &cobra.Command{
	Use:   "niche",
	Short: "Audit a YouTube channel + find niche opportunities",
	Long:  "Uses LLM + (optional) YouTube Data API v3 for live channel/video stats.",
	Example: `  ttb niche --channel @mkbhd
  ttb niche --channel @mrbeast --deep
  ttb niche --channel @mkbhd --no-youtube          # LLM-only mode`,
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
	nicheCmd.Flags().BoolVar(&flagNoYouTube, "no-youtube", false, "Force LLM-only mode (skip YouTube Data API)")
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
	ytClient := pickYouTubeClient()
	if flagNoYouTube {
		ytClient = nil
	}
	system, user := prompts.NicheWithData(ctx, ytClient, flagChannel, flagDeep, cfg.Lang)
	out, err := client.Complete(ctx, system, user)
	if err != nil {
		return err
	}
	if cfg.JSON {
		fmt.Println(out)
		return nil
	}
	header := fmt.Sprintf("Niche audit for %q", flagChannel)
	if ytClient != nil {
		header += " (with live YouTube data)"
	} else {
		header += " (LLM only)"
	}
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)))
	for _, line := range splitList(out) {
		fmt.Println("  - " + line)
	}
	return nil
}
