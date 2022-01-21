package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// re := `^.*?\"text\"\:\"(.*?)\"`
type Doc struct {
	DocId      string  `json:"review_id"`
	UserId     string  `json:"user_id"`
	BusinessId string  `json:"business_id"`
	Stars      float32 `json:"stars"`
	Useful     uint16  `json:"useful"`
	Funny      uint16  `json:"funny"`
	Cool       uint16  `json:"cool"`
	Text       string  `json:"text"`
	Date       string  `json:"date"`
}

func splitData(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := strings.Index(string(data), "}\n{"); i >= 0 {

		// Received partial JSON, tell caller to give us more data.
		firstRightBraceIndex := strings.Index(string(data), "}")
		if firstRightBraceIndex == -1 {
			return 0, nil, nil
		}

		secondLocation :=
			strings.Index(string(data[firstRightBraceIndex+1:]), "}")
		secondRightBraceIndex := secondLocation
		if secondLocation > 0 {
			secondRightBraceIndex = firstRightBraceIndex + secondLocation
		}

		if len(string(data)) > firstRightBraceIndex {
			if secondRightBraceIndex != -1 {
				return i + 1, data[:i+1], nil
			} else {
				return 0, nil, nil
			}
		}
		return i + 1, data[:i+1], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return
}

func LoadData(dataFile string) (<-chan string, <-chan struct{}) {
	done := make(chan struct{})
	in := make(chan string, 100)
	var doc Doc

	go func() {
		fmt.Println("In LoadData")
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
			if err := json.Unmarshal([]byte(line), &doc); err != nil {
				log.Fatalf("Error unmarshalling data: %s\n", line)
				os.Exit(-1)
			}

			in <- doc.Text
		}

		close(in)
		close(done)
	}()

	return in, done
}
