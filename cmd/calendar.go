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

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Plan a content calendar for a month",
	Long:  "Generates a topic-by-day plan with titles, angles, and post types.",
	Example: `  ttb calendar --month 7 --niche "AI tools" --frequency 2/week
  ttb calendar --month 2026-08 --niche "..." --frequency 3/week`,
	RunE: runCalendar,
}

var (
	flagMonth     string
	flagFrequency string
)

func init() {
	rootCmd.AddCommand(calendarCmd)
	calendarCmd.Flags().StringVar(&flagMonth, "month", "", "Target month: '7' or '2026-07' or 'YYYY-MM' (default: next month)")
	calendarCmd.Flags().StringVar(&flagNiche, "niche", "", "Niche/topic (required)")
	calendarCmd.Flags().StringVar(&flagFrequency, "frequency", "2/week", "Posting frequency: 1/week | 2/week | 3/week | daily")
	_ = calendarCmd.MarkFlagRequired("niche")
}

func runCalendar(cmd *cobra.Command, args []string) error {
	cfg, err := resolveConfig(cmd)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client, err := llm.New(cfg)
	if err != nil {
		return err
	}
	month := flagMonth
	if month == "" {
		month = "next month"
	}
	system, user := prompts.Calendar(month, flagNiche, flagFrequency, cfg.Lang)
	out, err := client.Complete(ctx, system, user)
	if err != nil {
		return err
	}
	if cfg.JSON {
		fmt.Println(out)
		return nil
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Content calendar for %q (%s, freq=%s):\n\n", month, flagNiche, flagFrequency)
	for _, line := range splitList(out) {
		fmt.Fprintln(tw, "  - "+line)
	}
	return tw.Flush()
}
