package main

import (
	"fmt"
	"os"

	"github.com/find-in-docs/sidecar/pkg/client"
	"github.com/find-in-docs/upload/pkg/config"
	"github.com/find-in-docs/upload/pkg/data"
	"github.com/find-in-docs/upload/pkg/transform"
	"github.com/spf13/viper"
)

func main() {

	sidecar, err := client.InitSidecar(tableName, nil)
	if err != nil {
		fmt.Printf("Error initializing sidecar: %v\n", err)
		os.Exit(-1)
	}

	config.LoadConfig()
	stopwords := data.LoadStopwords(viper.GetString("englishStopwordsFile"))

	wordsToInts := transform.WordsToInts(stopwords)
	for {
		v, ok := disk.LoadDoc()
		if !ok {
			break
		}
		wordInts, wordToInt = wordsToInts(v.Text)
		v.WordInts = wordInts
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
}
