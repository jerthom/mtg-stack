package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"mtgStack/app"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mtgStack",
	Short: "Generates an inventory of MTG cards from a given directory of images",
	Long: `mtgStack is a CLI that uses image text detection to identify Magic: the Gathering cards
			and produces a text file containing the counts and names for each image identified. `,
	RunE: func(cmd *cobra.Command, args []string) error {
		imgDir, _ := cmd.Flags().GetString("imagesDir")
		outputPath, _ := cmd.Flags().GetString("outputPath")

		cards, err := app.FetchCards(imgDir)
		if err != nil {
			return fmt.Errorf("unable to fetch cards: %v", err)
		}

		err = app.RecordCards(outputPath, cards)
		if err != nil {
			return fmt.Errorf("unable to record cards: %v", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("imagesDir", "i", "img", "Path to directory containing card images")
	rootCmd.Flags().StringP("outputPath", "o", "output.txt", "File path to write resulting card inventory")
}
