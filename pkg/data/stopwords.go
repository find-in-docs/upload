package data

import (
	"encoding/json"
	"fmt"
	"os"
)

func getJson(fn string, d interface{}) error {
	stopwordsFile, err := os.Open(fn)
	if err != nil {
		fmt.Printf("Error opening stopwords file: %v, %s", fn, err)
		return err
	}
	defer stopwordsFile.Close()

	jsonDecoder := json.NewDecoder(stopwordsFile)
	if err := jsonDecoder.Decode(d); err != nil {
		fmt.Printf("Error decoding file %s, %s\n", fn, d)
		return err
	}

	return nil
}

func Load(stopwordsFn string) []string {

	var stopwords []string

	if err := getJson(stopwordsFn, &stopwords); err != nil {
		fmt.Printf("Could not open stopwords file: %s\n", stopwordsFn)
		os.Exit(-1)
	}

	return stopwords
}
