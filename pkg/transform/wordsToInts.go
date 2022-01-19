package transform

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/samirgadkari/search/pkg/data"
)

type StringToStringFunc func(string) string
type StringToStringSliceFunc func(string, []string) []string
type StringSliceToIntSliceFunc func([]string, []int) []int
type IntSliceWriteToFileFunc func(string, int, []int)
type ProcFunc struct {
	ToLower       StringToStringFunc
	Replace       StringToStringFunc
	GetWords      StringToStringSliceFunc
	WordsToInts   StringSliceToIntSliceFunc
	WriteWordInts IntSliceWriteToFileFunc
}

const (
	maxWordsPerDoc = 1024
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

			if len(wordMatches) > maxWordsPerDoc {
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

func wordsToInts(stopwords []string, cfg Config) <-chan []int {

	proc := GenProcFunc(stopWords)

	in, done := data.LoadData(&cfg.DataFile)
	words := make([]string, maxWordsPerDoc)
	wordInts := make([]int, maxWordsPerDoc)
	processedReviews := make(chan []int)
	var line string
LOOP:
	for {
		select {
		case line = <-in:
			line = proc.Replace(line)
			line = proc.ToLower(line)
			words = proc.GetWords(line, words)
			wordInts = proc.WordsToInts(words, wordInts)
			fmt.Println(wordInts)
			processedReviews <- wordInts
			// type IntSliceWriteToFileFunc func(string, int, []int)
			// WriteWordInts IntSliceWriteToFileFunc
		case <-done:
			break LOOP
		}
	}

	return processedReviews
}
