package data

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

type InputDoc struct {
	InputDocId string  `json:"review_id"`
	UserId     string  `json:"user_id"`
	BusinessId string  `json:"business_id"`
	Stars      float32 `json:"stars"`
	Useful     uint16  `json:"useful"`
	Funny      uint16  `json:"funny"`
	Cool       uint16  `json:"cool"`
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
	Useful     uint16
	Funny      uint16
	Cool       uint16
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

func copyInputDocToDoc(inputDoc *InputDoc, doc *Doc) {
	doc.InputDocId = inputDoc.InputDocId
	doc.UserId = inputDoc.UserId
	doc.BusinessId = inputDoc.BusinessId
	doc.Stars = inputDoc.Stars
	doc.Useful = inputDoc.Useful
	doc.Funny = inputDoc.Funny
	doc.Cool = inputDoc.Cool
	doc.Text = inputDoc.Text
	doc.Date = inputDoc.Date
}

func LoadDocFn(dataFile string) func() (*Doc, bool) {
	done := make(chan struct{})
	in := make(chan Doc)
	var inputDoc InputDoc
	var doc *Doc = new(Doc)
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
			doc.DocId = DocumentId(docId)
			docId += 1
			in <- *doc
		}

		close(in)
		close(done)
	}()

	return func() (*Doc, bool) {
		line, ok := <-in
		// fmt.Printf("ok: %t, Got line: %d\n", ok, len(line))
		if !ok {
			return nil, ok
		}
		return &line, ok
	}
}
