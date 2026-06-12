package cmd
import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/llm"
	"github.com/arkydarmalik-coder/tube-trend-buddy/internal/prompts"
)

var thumbnailCmd = &cobra.Command{
	Use:   "thumbnail",
	Short: "Generate thumbnail concept (text overlay + visual prompt)",
	Long:  "Returns a JSON-ish concept you can feed to Midjourney / DALL-E / SD / Flux.",
	Example: `  ttb thumbnail --title "I Built an AI Agent" --mood shocked
  ttb thumbnail --title "..." --face 2 --colors "yellow,black,red"`,
	RunE: runThumbnail,
}

var (
	flagMood   string
	flagFace   int
	flagColors string
)

func init() {
	rootCmd.AddCommand(thumbnailCmd)
	thumbnailCmd.Flags().StringVar(&flagTitle, "title", "", "Video title (required)")
	thumbnailCmd.Flags().StringVar(&flagMood, "mood", "curious", "Mood: shocked | curious | excited | serious | funny")
	thumbnailCmd.Flags().IntVar(&flagFace, "face", 1, "Number of human faces in the thumbnail (0-3)")
	thumbnailCmd.Flags().StringVar(&flagColors, "colors", "", "Comma-separated color palette (e.g. 'yellow,black,red')")
	_ = thumbnailCmd.MarkFlagRequired("title")
}

func runThumbnail(cmd *cobra.Command, args []string) error {
	cfg, err := resolveConfig(cmd)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client, err := llm.New(cfg)
	if err != nil {
		return err
	}
	system, user := prompts.Thumbnail(flagTitle, flagMood, flagFace, flagColors, cfg.Lang)
	out, err := client.Complete(ctx, system, user)
	if err != nil {
		return err
	}
	if cfg.JSON {
		fmt.Println(out)
		return nil
	}
	fmt.Printf("Thumbnail concept for %q (mood=%s, faces=%d):\n\n", flagTitle, flagMood, flagFace)
	for _, line := range splitList(out) {
		fmt.Println("  - " + line)
	}
	return nil
}
