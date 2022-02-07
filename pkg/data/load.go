package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"
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

// Algorithm:
// 1. Read block of data. Use regex to split it into lines. Put lines in
// channel lines.
// 2. Read line. Convert to Document. Add DocId. Put Doc in channel docs.
//   From here we can fan out to process lines in parallel.
// 3. Read docs. Words to nums. Put (*Doc, []Num) in channel nums.
// 4. Read nums. Store each line in database.
// 5. Read nums. Store wordid-to-docid in database.
// 5. Store word-to-int mappings in database.

func readData() {

	dataFilename := viper.GetString("dataFile")
	lines := make(chan string)

	f, err := os.Open(dataFilename)
	if err != nil {
		fmt.Printf("Error opening file %s: err: %v\n", dataFilename, err)
		os.Exit(-1)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	chunkSize := 1024
	b := make([]byte, chunkSize)
	re := regexp.MustCompile(`\{.*?\}`)
	buf := make([]byte, 0)

	for _, err = r.Read(b); err != nil; {

		buf = append(buf, b...)
		for line := re.Find(buf); line != nil; {

			lines <- string(line)
		}
	}
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
