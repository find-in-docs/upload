package data

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	pb "github.com/find-in-docs/sidecar/protos/v1/messages"
)

type InputDoc struct {
	InputDocId string  `json:"review_id"`
	UserId     string  `json:"user_id"`
	BusinessId string  `json:"business_id"`
	Stars      float32 `json:"stars"`
	Useful     uint32  `json:"useful"`
	Funny      uint32  `json:"funny"`
	Cool       uint32  `json:"cool"`
	Text       string  `json:"text"`
	Date       string  `json:"date"`
}

const (
	MaxWordsPerDoc = 1024
)

type WordInt uint64
type DocumentId WordInt

type Doc struct {
	DocId      DocumentId
	WordInts   []WordInt
	InputDocId string
	UserId     string
	BusinessId string
	Stars      float32
	Useful     uint32
	Funny      uint32
	Cool       uint32
	Text       string
	Date       string
}

func copyInputDocToDoc(doc *pb.Doc, inputDoc *InputDoc) {
	doc.InputDocId = strings.ToValidUTF8(inputDoc.InputDocId, " ")
	doc.UserId = strings.ToValidUTF8(inputDoc.UserId, " ")
	doc.BusinessId = strings.ToValidUTF8(inputDoc.BusinessId, " ")
	doc.Stars = inputDoc.Stars
	doc.Useful = inputDoc.Useful
	doc.Funny = inputDoc.Funny
	doc.Cool = inputDoc.Cool
	doc.Text = strings.ToValidUTF8(inputDoc.Text, " ")
	doc.Date = strings.ToValidUTF8(inputDoc.Date, " ")
}

func LoadDocFn(dataFile string) chan *pb.Doc {
	in := make(chan *pb.Doc)
	var inputDoc InputDoc

	go func() {
		f, err := os.Open(dataFile)
		if err != nil {
			log.Fatalf("Error opening file %s", dataFile)
			os.Exit(-1)
		}

		defer f.Close()

		s := bufio.NewScanner(f)

		for s.Scan() {
			line := s.Text()
			if err := json.Unmarshal([]byte(line), &inputDoc); err != nil {
				log.Fatalf("Error unmarshalling data: %s\n", line)
				os.Exit(-1)
			}

			doc := new(pb.Doc)
			copyInputDocToDoc(doc, &inputDoc)
			in <- doc
		}

		close(in)
	}()

	return in
}
