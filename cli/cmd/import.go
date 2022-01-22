/*
Copyright © 2022 Samir Gadkari

*/
package cmd

import (
	"github.com/samirgadkari/search/pkg/config"
	"github.com/samirgadkari/search/pkg/transform"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import documents from the given file",
	Long: `Specify a file with documents in JSON object form.
If it is a list of documents, don't include the [] list specifiers. ex:
{"review_id": 1, "text": "User review for ID 1"}
{"review_id": 2, "text": "User review for ID 2"}`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		stopWords := config.LoadStopwords(cfg)

		var dataFile string
		if len(args) == 0 {
			dataFile = cfg.DataFile
		} else {
			dataFile = args[0]
		}

		transform.WordsToInts(stopWords, dataFile,
			cfg.OutputDir, cfg.WordIntsFile)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
