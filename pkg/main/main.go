package main

import (
	"context"
	"fmt"
	"os"

	"github.com/find-in-docs/sidecar/pkg/client"
	pb "github.com/find-in-docs/sidecar/protos/v1/messages"
	"github.com/find-in-docs/upload/pkg/config"
	"github.com/find-in-docs/upload/pkg/data"
	"github.com/spf13/viper"
	proto "google.golang.org/protobuf/proto"
)

const (
	chunkSize             = 5
	allTopicsRecvChanSize = 32
)

func main() {

	config.Load()

	fmt.Printf("sidecarServiceAddr: %s\n", viper.GetString("sidecarServiceAddr"))
	sidecar, err := client.InitSidecar(viper.GetString("serviceName"), nil)
	if err != nil {
		fmt.Printf("Error initializing sidecar: %v\n", err)
		os.Exit(-1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	docsCh := data.LoadDocFn(viper.GetString("dataFile"))

	documents := new(pb.Documents)
	docs := make([]*pb.Doc, chunkSize)
	documents.Docs = docs
	var count int
	var numOutput int
	for doc := range docsCh {

		docs[count] = doc
		if count == chunkSize-1 {

			fmt.Printf("len(docs): %d\n\n", len(docs))
			b, err := proto.Marshal(documents)
			if err != nil {
				fmt.Printf("Error encoding document: %v\n", err)
				return
			}

			// Publish data to message queue
			err = sidecar.PubJS(ctx, "search.doc.import.v1", "uploadWorkQueue", b)
			if err != nil {
				fmt.Printf("Error publishing message.\n\terr: %v\n", err)
			}

			count = 0
			numOutput++
			if numOutput == 2 {
				break
			}
		} else {

			count++
		}
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

	select {} // This will wait forever
}
