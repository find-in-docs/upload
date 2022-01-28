/*
Copyright © 2022 Samir Gadkari

*/
package cmd

import (
	"fmt"
	"os"
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

		config.LoadConfig()

		stopwords := data.LoadStopwords(viper.GetString("englishStopwordsFile"))

		switch viper.GetString("output.type") {
		case config.File.String():
			outputDir = filepath.Dir(viper.GetString("output.location"))
			wordIntsFile = filepath.Base(viper.GetString("output.location"))

			disk := data.DiskSetup(outputDir, wordIntsFile)

			var wordInts []data.WordInt
			var wordToInt map[string]data.WordInt
			var intToWord map[data.WordInt]string

			wordsToInts := transform.WordsToInts(stopwords)
			for {
				v, ok := disk.LoadDoc()
				if !ok {
					break
				}
				wordInts, wordToInt, intToWord = wordsToInts(v.Text)
				disk.StoreData(v, wordInts)
			}

			disk.WriteWordIntMappings(wordToInt, intToWord)
			disk.Close()

		case config.Database.String():
			fmt.Println("Database support in progress")

			outputDir = filepath.Dir(viper.GetString("output.location"))
			wordIntsFile = filepath.Base(viper.GetString("output.location"))

			disk := data.DiskSetup(outputDir, wordIntsFile)

			db := data.DBSetup()

			var wordInts []data.WordInt
			var wordToInt map[string]data.WordInt
			var intToWord map[data.WordInt]string
			tableName := "doc"

			db.OpenConnection()
			db.CreateTable(tableName)

			wordsToInts := transform.WordsToInts(stopwords)
			for {
				v, ok := disk.LoadDoc()
				if !ok {
					break
				}
				wordInts, wordToInt, intToWord = wordsToInts(v.Text)
				v.WordInts = wordInts
				db.StoreData(v, tableName, wordInts)
			}

			switch viper.GetString("output.type") {
			case "file":
				disk.WriteWordIntMappings(wordToInt, intToWord)
			case "database":
				db.StoreWordIntMappings("wordtoint", wordToInt, "inttoword", intToWord)
			default:
				fmt.Printf("Incorrect output.type in config file\n")
				os.Exit(-1)
			}

			db.CloseConnection()
			disk.Close()
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
