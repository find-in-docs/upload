package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

func process(doc *Doc, r *strings.Replacer, re *regexp.Regexp) *[]string {
	s := strings.ToLower(doc.Text)
	s = r.Replace(s)

	matches := re.FindAllStringSubmatch(s, -1)
	if matches != nil {
		wordMatches := make([]string, len(matches))
		matchNum := 0
		for _, match := range matches {
			if len(match) != 2 {
				fmt.Printf("len(match) is %d which is invalid\n", len(match))
				os.Exit(-1)
			}
			wordMatches[matchNum] = match[1]
			matchNum += 1
		}

		return &wordMatches
	} else {
		return nil
	}
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
	docNum := 0

	// isn't -> isnt
	// solar_eclipse -> solar eclipse
	// We also remove _ because it is a part of \w rexexp,
	// and we don't want the whole "solar_eclipse" as a word.
	// We're still ok with 0-9 being part of the word,
	// although maybe that may change after some experience.
	r := strings.NewReplacer("'", "", "_", " ")

	re := regexp.MustCompile(`(\w+)`)
	for yamlDecoder.More() {
		if err = yamlDecoder.Decode(doc); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%10v doc: %#v\n\n", docNum, doc)
		_ = process(doc, r, re)
		docNum += 1
	}

	fmt.Printf("Number of docs: %d", docNum)
}
