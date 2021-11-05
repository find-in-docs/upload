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

type Config struct {
	RawDocumentsFn string `yaml:"raw-documents-fn"`
	OutputDirFn    string `yaml:"output-dir-fn"`
	StopwordsFn    string `yaml:"stopwords-fn"`
}

type Stopwords struct {
	English []string `json:"english_stopwords"`
}

type StringToStringFunc func(string) string
type StringToStringSliceFunc func(string, []string) []string
type StringSliceToIntSliceFunc func([]string, []int) []int
type ProcFunc struct {
	ToLower     StringToStringFunc
	Replace     StringToStringFunc
	GetWords    StringToStringSliceFunc
	WordsToInts StringSliceToIntSliceFunc
}

const (
	MAX_NUM_WORDS_PER_DOC = 1024
)

func GenProcFunc(stopwords []string) *ProcFunc {

	var procFunc ProcFunc

	// This character replacer replaces:
	// an apostrophe ' with empty string,
	// a underscore _ with space,
	// and all other punctuation marks with space.
	// We're also going to remove all numbers.
	r := strings.NewReplacer("'", "", "_", " ", "!", " ",
		"!", " ", "\"", " ", "#", " ", "$", " ",
		"%", " ", "&", " ",
		"(", " ", ")", " ", "*", " ", "+", " ",
		",", " ", "-", " ", ".", " ", "/", " ",
		":", " ", ";", " ", "<", " ", "=", " ",
		">", " ", "?", " ", "@", " ", "[", " ",
		"\\", " ", "]", " ", "^", " ",
		"_", " ", "`", " ", "{", " ", "|", " ",
		"}", " ", "~", " ",
		"0", "", "1", "", "2", "", "3", "", "4", "",
		"5", "", "6", "", "7", "", "8", "", "9", "")

	procFunc.Replace = func(s string) string {
		return r.Replace(s)
	}

	procFunc.ToLower = func(s string) string {
		return strings.ToLower(s)
	}

	re := regexp.MustCompile(`(\w+)`)
	procFunc.GetWords = func(s string, wordMatches []string) []string {
		matches := re.FindAllStringSubmatch(s, -1)
		if matches != nil {
			wordMatches = wordMatches[:0]
			for _, match := range matches {
				if len(match) != 2 {
					fmt.Printf("len(match) is %d which is invalid\n", len(match))
					os.Exit(-1)
				}
				wordMatches = append(wordMatches, match[1])
			}

			if len(wordMatches) > MAX_NUM_WORDS_PER_DOC {
				fmt.Printf("Too many words in doc (%d) !!", len(wordMatches))
				os.Exit(-1)
			}

			return wordMatches
		} else {
			return nil
		}
	}

	wordToInt := make(map[string]int)
	IntToWord := make(map[int]string)
	wordNum := 0
	procFunc.WordsToInts = func(words []string, wordInts []int) []int {

		wordInts = wordInts[:0]
		for _, word := range words {
			if _, ok := wordToInt[word]; ok == false {
				wordToInt[word] = wordNum
				IntToWord[wordNum] = word
				wordInts = append(wordInts, wordNum)
				wordNum += 1
			} else {
				wordInts = append(wordInts, wordToInt[word])
			}

		}

		return wordInts
	}

	return &procFunc
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

func main() {

	workingDir, _ := os.Getwd()
	fmt.Printf("pwd: %s\n", workingDir)

	absConfigPath, _ := filepath.Abs("../../../config.yaml")
	configFile, err := os.Open(absConfigPath)
	if err != nil {
		fmt.Printf("Error opening file: %s, %s", absConfigPath, err)
		os.Exit(-1)
	}
	defer configFile.Close()

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
	fmt.Printf("config: %#v\n", config)

	var stopwords Stopwords
	if getJson(&config.StopwordsFn, &stopwords) != nil {
		fmt.Printf("Error getting stopwords from file\n")
		os.Exit(1)
	}
	fmt.Printf("stopwords.English: %#v\n", stopwords.English)

	filename := config.RawDocumentsFn
	fmt.Printf("Data filename: %s\n", filename)
	dataFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening data file: %s, %s", filename, err)
	}
	defer dataFile.Close()

	jsonDecoder := json.NewDecoder(dataFile)
	doc := &Doc{}
	docNum := 0
	var s string
	words := make([]string, MAX_NUM_WORDS_PER_DOC)
	wordInts := make([]int, MAX_NUM_WORDS_PER_DOC)

	procFunc := GenProcFunc(stopwords.English)
	for jsonDecoder.More() {
		if err = jsonDecoder.Decode(doc); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%10v doc: %#v\n\n", docNum, doc)

		s = procFunc.ToLower(doc.Text)
		s = procFunc.Replace(s)
		words = procFunc.GetWords(s, words)
		fmt.Printf("words: %#v\n", words)
		wordInts = procFunc.WordsToInts(words, wordInts)
		fmt.Printf("wordInts: %#v\n", wordInts)
		docNum += 1
	}

	fmt.Printf("Number of docs: %d", docNum)
}
