package data

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

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

	dataString := string(data)
	if i := strings.Index(dataString, "}"); i >= 0 {
		return i + 2, data[:i+2], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func LoadData(dataFile string) (<-chan string, <-chan struct{}) {
	done := make(chan struct{})
	in := make(chan string)
	var doc Doc

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
