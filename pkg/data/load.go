package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

func getJson(fn *string, d interface{}) error {
	fmt.Printf("Decoding JSON file: %s\n", *fn)
	stopwordsFile, err := os.Open(*fn)
	if err != nil {
		fmt.Printf("Error opening stopwords file: %s, %s", fn, err)
	}
	defer stopwordsFile.Close()

	jsonDecoder := json.NewDecoder(stopwordsFile)
	if err := jsonDecoder.Decode(d); err != nil {
		fmt.Printf("Error decoding file %s, %s\n", *fn, d)
		os.Exit(-1)
	}

	return nil
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
		s.Split(bufio.ScanLines)

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
