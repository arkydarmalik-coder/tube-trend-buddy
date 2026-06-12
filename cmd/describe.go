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

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Write SEO description + hashtags + CTA for a video",
	Long:  "Crafts a search-friendly description with optional timestamps and pinned comment CTA.",
	Example: `  ttb describe --title "How to Learn AI" --tags "AI, learning"
  ttb describe --title "..." --cta --timestamps "00:00 intro, 02:30 basics"`,
	RunE: runDescribe,
}

var (
	flagCTA        bool
	flagTimestamps string
)

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringVar(&flagTitle, "title", "", "Video title (required)")
	describeCmd.Flags().StringVar(&flagDescription, "description", "", "Optional short summary / bullet points")
	describeCmd.Flags().StringVar(&flagTimestamps, "timestamps", "", "Optional comma-separated chapter list (e.g. '00:00 intro,02:30 basics')")
	describeCmd.Flags().BoolVar(&flagCTA, "cta", true, "Include a call-to-action block")
	_ = describeCmd.MarkFlagRequired("title")
}

func runDescribe(cmd *cobra.Command, args []string) error {
	cfg, err := resolveConfig(cmd)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client, err := llm.New(cfg)
	if err != nil {
		return err
	}
	system, user := prompts.Describe(flagTitle, flagDescription, flagTimestamps, flagCTA, cfg.Lang)
	out, err := client.Complete(ctx, system, user)
	if err != nil {
		return err
	}
	if cfg.JSON {
		fmt.Println(out)
		return nil
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Description for %q:\n\n", flagTitle)
	fmt.Fprintln(tw, out)
	return tw.Flush()
}
