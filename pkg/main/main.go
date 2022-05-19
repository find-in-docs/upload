package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/find-in-docs/sidecar/pkg/client"
	pb "github.com/find-in-docs/sidecar/protos/v1/messages"
	"github.com/find-in-docs/upload/pkg/config"
	"github.com/find-in-docs/upload/pkg/data"
	"github.com/spf13/viper"
	proto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
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

	topic := "search.doc.import.response.v1"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	msgStrRegex := regexp.MustCompile(`\\+?\"|\\+?n|\\+?t`)

	err = sidecar.ProcessSubMsgs(ctx, topic,
		allTopicsRecvChanSize, func(m *pb.SubTopicResponse) {

			msg := fmt.Sprintf("Received from sidecar:\n\t%s", m.String())
			msg2 := msgStrRegex.ReplaceAllString(msg, "")
			fmt.Printf("%s\n", msg2)

			// Process incoming message
		})

	if err != nil {
		fmt.Printf("Error processing subscription messages:\n\ttopic: %s\n\terr: %v\n",
			topic, err)
	}

	docsFn := data.LoadDocFn(viper.GetString("dataFile"))

	var retryNum uint32 = 1
	retryDelayDuration, err := time.ParseDuration("200ms")
	if err != nil {
		fmt.Printf("Error creating Golang time duration.\nerr: %v\n", err)
		os.Exit(-1)
	}
	retryDelay := durationpb.New(retryDelayDuration)

	var doc *pb.Doc
	documents := new(pb.Documents)
	docs := make([]*pb.Doc, chunkSize)
	documents.Docs = docs
	var ok bool
	var count int
	for doc, ok = docsFn(); ok; {

		docs[count] = doc
		if count%(chunkSize-1) == 0 {

			b, err := proto.Marshal(documents)
			if err != nil {
				fmt.Printf("Error encoding document: %v\n", err)
				return
			}

			// Publish data to message queue
			err = sidecar.Pub(ctx, "search.doc.import.v1", b,
				&pb.RetryBehavior{
					RetryNum:   &retryNum,
					RetryDelay: retryDelay,
				},
			)
			if err != nil {
				fmt.Printf("Error publishing message.\n\terr: %v\n", err)
			}

			count = 0
			break
		}

		count++
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
}
