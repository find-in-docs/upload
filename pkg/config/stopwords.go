package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func getJson(fn *string, d interface{}) error {
	fmt.Printf("Decoding JSON file: %s\n", *fn)
	stopwordsFile, err := os.Open(*fn)
	if err != nil {
		fmt.Printf("Error opening stopwords file: %v, %s", fn, err)
		return err
	}
	defer stopwordsFile.Close()

	jsonDecoder := json.NewDecoder(stopwordsFile)
	if err := jsonDecoder.Decode(d); err != nil {
		fmt.Printf("Error decoding file %s, %s\n", *fn, d)
		return err
	}

	return nil
}

func LoadStopwords(cfg *Config) []string {

	var stopwords []string

	if err := getJson(&cfg.EnglishStopwordsFile, &stopwords); err != nil {
		os.Exit(-1)
	}

	return stopwords
}
