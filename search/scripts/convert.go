package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Doc struct {
	DocId      string  `json: "review_id"`
	UserId     string  `json: "user_id"`
	BusinessId string  `json: "business_id"`
	Stars      float32 `json: "stars"`
	Useful     uint16  `json: "useful"`
	Funny      uint16  `json: "funny"`
	Cool       uint16  `json: "cool"`
	Text       string  `json: "text"`
	Date       string  `json: "date"`
}

type Config struct {
	RawDocumentsFn string `yaml:"raw-documents-fn"`
	OutputDirFn    string `yaml:"output-dir-fn"`
	StopwordsFn    string `yaml:"stopwords-fn"`
}

func main() {

	workingDir, _ := os.Getwd()
	fmt.Printf("pwd: %s\n", workingDir)

	absConfigPath, _ := filepath.Abs("../../../config.yaml")
	configFile, err := os.Open(absConfigPath)
	if err != nil {
		fmt.Printf("Error opening file: %s, %s", absConfigPath, err)
		os.Exit(-1)
	}

	var config Config
	decoder := yaml.NewDecoder(configFile)
	if decoder == nil {
		fmt.Printf("Error getting new decoder for file: %s", absConfigPath)
		os.Exit(-1)
	}

	if err = decoder.Decode(&config); err != nil {
		fmt.Printf("Error decoding file %s", absConfigPath)
		os.Exit(-1)
	}
	fmt.Printf("config: %#v", config)

	filename := config.RawDocumentsFn
	fmt.Printf("Data filename: %s\n", filename)
	dataFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening data file: %s, %s", filename, err)
	}
	defer dataFile.Close()

	yamlDecoder := json.NewDecoder(dataFile)
	if yamlDecoder == nil {
		fmt.Printf("Error getting new decoder for file: %s", filename)
		os.Exit(-1)
	}

	doc := &Doc{}
	i := 0
	for yamlDecoder.More() {
		if err = yamlDecoder.Decode(doc); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%10v doc: %#v\n\n", i, doc)
		i += 1
	}

	fmt.Printf("Number of docs: %d", i)
}
