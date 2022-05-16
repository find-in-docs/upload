package transform

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/find-in-docs/search/pkg/data"
)

type ProcFunc struct {
	ToLower            func(string) string
	Replace            func(string) string
	GetWords           func(string, []string) []string
	WordsToInts        func([]string, []data.WordInt) []data.WordInt
	WriteWordInts      func(string, int, []data.WordInt)
	GetWordIntMappings func() map[string]data.WordInt
	RemoveStopwords    func([]string) []string
}

const (
	wordToIntFilename = "wordToInt.txt"
	intToWordFilename = "intToWord.txt"
)

func replacer(s string) string {
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

	return r.Replace(s)
}

func getWordsFn() func(string, []string) []string {

	re := regexp.MustCompile(`(\w+)`)

	return func(s string, wordMatches []string) []string {
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

			if len(wordMatches) > data.MaxWordsPerDoc {
				fmt.Printf("Too many words in doc (%d) !!", len(wordMatches))
				os.Exit(-1)
			}

			return wordMatches
		} else {
			return nil
		}
	}
}

func removeStopwordsFn(stopwords []string) func([]string) []string {

	swMap := make(map[string]struct{}, len(stopwords))
	for _, v := range stopwords {
		swMap[v] = struct{}{}
	}

	return func(words []string) []string {

		var result []string
		for _, word := range words {
			_, ok := swMap[word]
			if !ok {
				result = append(result, word)
			}
		}

		return result
	}
}

func wordToIntsFns() (func([]string, []data.WordInt) []data.WordInt,
	func() map[string]data.WordInt) {

	wordToInt := make(map[string]data.WordInt)
	var wordNum data.WordInt = 0
	return func(words []string, wordInts []data.WordInt) []data.WordInt {

			wordInts = wordInts[:0]
			for _, word := range words {
				if _, ok := wordToInt[word]; ok == false {
					wordToInt[word] = wordNum
					wordInts = append(wordInts, wordNum)
					wordNum += 1
				} else {
					wordInts = append(wordInts, wordToInt[word])
				}

			}

			return wordInts
		}, func() map[string]data.WordInt {
			return wordToInt
		}
}

func GenProcFunc(stopwords []string) *ProcFunc {

	var procFunc ProcFunc

	procFunc.Replace = replacer

	procFunc.ToLower = func(s string) string {
		return strings.ToLower(s)
	}

	procFunc.GetWords = getWordsFn()
	procFunc.RemoveStopwords = removeStopwordsFn(stopwords)
	procFunc.WordsToInts, procFunc.GetWordIntMappings = wordToIntsFns()

	return &procFunc
}

func WordsToInts(stopwords []string) func(string) ([]data.WordInt, map[string]data.WordInt) {

	proc := GenProcFunc(stopwords)

	words := make([]string, data.MaxWordsPerDoc)
	wordInts := make([]data.WordInt, data.MaxWordsPerDoc)

	return func(line string) ([]data.WordInt, map[string]data.WordInt) {

		line = proc.Replace(line)
		line = proc.ToLower(line)
		words = proc.GetWords(line, words)
		words = proc.RemoveStopwords(words)
		wordInts = proc.WordsToInts(words, wordInts)
		wordToInt := proc.GetWordIntMappings()

		return wordInts, wordToInt
	}
}

func WordToIntSwitchKV(wordToInt map[string]data.WordInt) *map[data.WordInt]string {

	intToWord := make(map[data.WordInt]string)
	for k, v := range wordToInt {
		intToWord[v] = k
	}

	return &intToWord
}
