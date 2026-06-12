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

var titleCmd = &cobra.Command{
	Use:   "title",
	Short: "Generate CTR-optimized YouTube titles for a niche",
	Long: "Analyzes top-performing videos in a niche and generates click-worthy titles.",
	Example: `  ttb title --niche "AI tools" --count 10
  ttb title --niche "resep Nusantara" --lang id --count 5`,
	RunE: runTitle,
}

var (
	flagNiche    string
	flagAudience string
)

func init() {
	rootCmd.AddCommand(titleCmd)
	titleCmd.Flags().StringVar(&flagNiche, "niche", "", "Niche/topic keyword (required)")
	titleCmd.Flags().StringVar(&flagAudience, "audience", "general", "Target audience description")
	_ = titleCmd.MarkFlagRequired("niche")
}

func runTitle(cmd *cobra.Command, args []string) error {
	shared.applyRootFlags(cmd)
	cfg, err := resolveConfig(cmd)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client, err := llm.New(cfg)
	if err != nil {
		return err
	}
	system, user := prompts.Title(flagNiche, flagAudience, cfg.Lang, cfg.Count)
	out, err := client.Complete(ctx, system, user)
	if err != nil {
		return err
	}
	if cfg.JSON {
		fmt.Println(out)
		return nil
	}
	fmt.Printf("Titles for niche %q (%d ideas, lang=%s):\n\n", flagNiche, cfg.Count, cfg.Lang)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, line := range splitList(out) {
		fmt.Fprintln(tw, "  - "+line)
	}
	return tw.Flush()
}
