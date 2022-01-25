/*
Copyright © 2022 Samir Gadkari

*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/samirgadkari/search/pkg/config"
	"github.com/samirgadkari/search/pkg/data"
	"github.com/samirgadkari/search/pkg/transform"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		var outputDir = ""
		var wordIntsFile = ""
		var dataFile = ""

		config.LoadConfig()

		if len(args) == 0 {
			dataFile = viper.GetString("output.location")
		} else {
			dataFile = args[0]
		}

		switch viper.GetString("output.type") {
		case config.File.String():
			outputDir = filepath.Dir(viper.GetString("output.location"))
			wordIntsFile = filepath.Base(viper.GetString("output.location"))
			fmt.Printf("outputDir: %s\nwordIntsFile: %s\n", outputDir, wordIntsFile)

			storeData, closeData := data.StoreDataOnDisk(outputDir, wordIntsFile)

			transform.WordsToInts(config.LoadStopwords,
				data.LoadDocFn(dataFile),
				storeData,
				closeData,
				outputDir)
		case config.Database.String():
			fmt.Println("Database support is progress\n")
		}
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
