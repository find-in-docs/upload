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

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	config.Load()

	fmt.Printf("sidecarServiceAddr: %s\n", viper.GetString("sidecarServiceAddr"))
	sidecar, err := client.InitSidecar(viper.GetString("serviceName"), nil)
	if err != nil {
		fmt.Printf("Error initializing sidecar: %v\n", err)
		os.Exit(-1)
	}

	wg.Add(1)

	subject := viper.GetString("nats.jetstream.subject")
	workQueue := viper.GetString("nats.jetstream.name")

  fmt.Printf("Adding subject:%s, workQueue:%s\n", subject, workQueue)
	if err = sidecar.AddJS(ctx, subject, workQueue); err != nil {

		fmt.Printf("Error adding stream: %v\n", err)
		os.Exit(-1)
	}

	docsCh := data.LoadDocFn(viper.GetString("dataFile"))

	if err != nil {
		fmt.Printf("Error creating stream. err: %v\n", err)
		os.Exit(-1)
	}

	if err = sidecar.UploadDocs(wg, docsCh); err != nil {
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

	// Cancel any goroutines which are still running for the upload
	cancel()

	if err = sidecar.UnsubJS(ctx, subject, workQueue); err != nil {

		fmt.Printf("Error unsubscribing from stream: %v\n", err)
		os.Exit(-1)
	}

	wg.Wait()
}
