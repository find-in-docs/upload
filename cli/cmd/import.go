/*
Copyright © 2022 Samir Gadkari

*/
package cmd

import (
	"fmt"

	"github.com/find-in-docs/upload/pkg/config"
	"github.com/find-in-docs/upload/pkg/data"
	"github.com/find-in-docs/upload/pkg/transform"
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

		config.LoadConfig()

		stopwords := data.LoadStopwords(viper.GetString("englishStopwordsFile"))

		disk := data.DiskSetup()

		switch viper.GetString("output.type") {
		case config.File.String():

			var wordInts []data.WordInt
			var wordToInt map[string]data.WordInt

			wordsToInts := transform.WordsToInts(stopwords)
			for {
				v, ok := disk.LoadDoc()
				if !ok {
					break
				}
				wordInts, wordToInt = wordsToInts(v.Text)
				disk.StoreData(v, wordInts)
			}

			disk.WriteWordIntMappings(wordToInt)
			disk.Close()

		case config.Database.String():

			dbFunc := data.DBSetup()

			var wordInts []data.WordInt
			var wordToInt map[string]data.WordInt
			tableName := "doc"

			if err := dbFunc.OpenConnection(); err != nil {
				break
			}
			if err := dbFunc.CreateTable(tableName); err != nil {
				break
			}

			wordsToInts := transform.WordsToInts(stopwords)
			for {
				v, ok := disk.LoadDoc()
				if !ok {
					break
				}
				wordInts, wordToInt = wordsToInts(v.Text)
				v.WordInts = wordInts
				if err := dbFunc.StoreData(v, tableName, wordInts); err != nil {
					break
				}
			}

			if err := dbFunc.StoreWordIntMappings("word_to_int", wordToInt); err != nil {
				break
			}

			fmt.Printf("Loading docs\n")
			inputDocs, err := data.LoadDocs()
			if err != nil {
				break
			}

			fmt.Printf("Transforming WordToDocs\n")
			if err := transform.WordToDocs(inputDocs, dbFunc.StoreWordToDocMappings); err != nil {
				break
			}

			if err := dbFunc.CloseConnection(); err != nil {
				break
			}
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
