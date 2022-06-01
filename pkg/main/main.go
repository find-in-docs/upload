package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/find-in-docs/sidecar/pkg/client"
	"github.com/find-in-docs/upload/pkg/config"
	"github.com/find-in-docs/upload/pkg/data"
	"github.com/spf13/viper"
)

func main() {

	var wg sync.WaitGroup

	config.Load()

	fmt.Printf("sidecarServiceAddr: %s\n", viper.GetString("sidecarServiceAddr"))
	sidecar, err := client.InitSidecar(viper.GetString("serviceName"), nil)
	if err != nil {
		fmt.Printf("Error initializing sidecar: %v\n", err)
		os.Exit(-1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)

	if err = sidecar.AddJS(ctx, viper.GetString("nats.jetstream.subject"),
		viper.GetString("nats.jetstream.name")); err != nil {

		fmt.Printf("Error adding stream: %v\n", err)
		os.Exit(-1)
	}

	docsCh := data.LoadDocFn(viper.GetString("dataFile"))

	if err != nil {
		fmt.Printf("Error creating stream. err: %v\n", err)
		os.Exit(-1)
	}

	if err = sidecar.UploadDocs(ctx, docsCh); err != nil {
		fmt.Printf("Error uploading docs. err: %v\n", err)
		os.Exit(-1)
	}

	/*

		We will have to do this in the documents service. For the upload service, we just
		send a chunk of documents to NATS, and allow the documents service to download
		and process each chunk in this fashion:

		stopwords := data.LoadStopwords(viper.GetString("englishStopwordsFile"))
		wordsToIntsFn := transform.WordsToInts(stopwords)
		var doc *data.Doc
		var ok bool
		var wordInts []data.WordInt
		var wordToInt map[string]data.WordInt
		for doc, ok = docs(); ok; {
			wordInts, wordToInt = wordsToIntsFn(doc.Text)
			doc.WordInts = wordInts
		}
	*/

	wg.Wait()
}
