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

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Generate SEO tags (rank + traffic) for a video",
	Long:  "Produces a balanced mix of high-volume, mid-tail and long-tail tags.",
	Example: `  ttb tags --title "How to Learn AI in 30 Days" --count 30
  ttb tags --title "..." --description "..." --count 20`,
	RunE: runTags,
}

var (
	flagTitle       string
	flagDescription string
)

func init() {
	rootCmd.AddCommand(tagsCmd)
	tagsCmd.Flags().StringVar(&flagTitle, "title", "", "Video title (required)")
	tagsCmd.Flags().StringVar(&flagDescription, "description", "", "Optional video description for context")
	_ = tagsCmd.MarkFlagRequired("title")
}

func runTags(cmd *cobra.Command, args []string) error {
	cfg, err := resolveConfig(cmd)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client, err := llm.New(cfg)
	if err != nil {
		return err
	}
	system, user := prompts.Tags(flagTitle, flagDescription, cfg.Lang, cfg.Count)
	out, err := client.Complete(ctx, system, user)
	if err != nil {
		return err
	}
	if cfg.JSON {
		fmt.Println(out)
		return nil
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "  Tags for %q (%d tags):\n\n", flagTitle, cfg.Count)
	for _, line := range splitList(out) {
		fmt.Fprintln(tw, "  # "+line)
	}
	return tw.Flush()
}
