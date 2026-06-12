package cmd

import (
	"github.com/spf13/cobra"

	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/ui"
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Launch the interactive TUI wizard",
	Long: `Opens a full-screen terminal UI to pick a feature, fill in inputs, and
view the LLM result. Built on Bubble Tea + Lip Gloss.`,
	Example: `  ttb ui`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ui.Run()
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
