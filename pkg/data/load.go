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

func splitData(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	dataString := string(data)
	if i := strings.Index(dataString, "}"); i >= 0 {
		return i + 2, data[:i+2], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func copyInputDocToDoc(inputDoc *InputDoc, doc *pb.Doc) {
	doc.InputDocId = strings.ToValidUTF8(inputDoc.InputDocId, " ")
	doc.UserId = strings.ToValidUTF8(inputDoc.UserId, " ")
	// doc.BusinessId = strings.ToValidUTF8(inputDoc.BusinessId, " ")
	doc.Stars = inputDoc.Stars
	doc.Useful = inputDoc.Useful
	doc.Funny = inputDoc.Funny
	doc.Cool = inputDoc.Cool
	doc.Text = strings.ToValidUTF8(inputDoc.Text, " ")
	doc.Date = strings.ToValidUTF8(inputDoc.Date, " ")
}

func LoadDocFn(dataFile string) func() (*pb.Doc, bool) {
	in := make(chan *pb.Doc)
	var inputDoc InputDoc
	var doc *pb.Doc = new(pb.Doc)
	var docId WordInt = 0

	go func() {
		f, err := os.Open(dataFile)
		defer f.Close()

		if err != nil {
			log.Fatalf("Error opening file %s", dataFile)
			os.Exit(-1)
		}

		s := bufio.NewScanner(f)

		s.Split(splitData)

		for s.Scan() {
			line := s.Text()
			if err := json.Unmarshal([]byte(line), &inputDoc); err != nil {
				log.Fatalf("Error unmarshalling data: %s\n", line)
				os.Exit(-1)
			}

			copyInputDocToDoc(&inputDoc, doc)
			doc.DocId = uint64(docId)
			docId += 1
			in <- doc
		}

		close(in)
	}()

	return func() (*pb.Doc, bool) {
		line, ok := <-in
		// fmt.Printf("ok: %t, Got line: %d\n", ok, len(line))
		if !ok {
			return nil, ok
		}
		return line, ok
	}
}
