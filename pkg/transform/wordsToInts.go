package transform

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/samirgadkari/search/pkg/data"
)

type StringToStringFunc func(string) string
type StringToStringSliceFunc func(string, []string) []string
type StringSliceToIntSliceFunc func([]string, []int) []int
type IntSliceWriteToFileFunc func(string, int, []int)
type ProcFunc struct {
	ToLower              StringToStringFunc
	Replace              StringToStringFunc
	GetWords             StringToStringSliceFunc
	WordsToInts          StringSliceToIntSliceFunc
	WriteWordInts        IntSliceWriteToFileFunc
	WriteWordIntMappings func(string)
	RemoveStopwords      func([]string) []string
}

const (
	maxWordsPerDoc    = 1024
	wordToIntFilename = "wordToInt.txt"
	intToWordFilename = "intToWord.txt"
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

	swMap := make(map[string]struct{}, len(stopwords))
	for _, v := range stopwords {
		swMap[v] = struct{}{}
	}

	procFunc.RemoveStopwords = func(words []string) []string {

		var result []string
		for _, word := range words {
			_, ok := swMap[word]
			if !ok {
				result = append(result, word)
			}
		}

		return result
	}

	wordToInt := make(map[string]int)
	intToWord := make(map[int]string)
	wordNum := 0
	procFunc.WordsToInts = func(words []string, wordInts []int) []int {

		wordInts = wordInts[:0]
		for _, word := range words {
			if _, ok := wordToInt[word]; ok == false {
				wordToInt[word] = wordNum
				intToWord[wordNum] = word
				wordInts = append(wordInts, wordNum)
				wordNum += 1
			} else {
				wordInts = append(wordInts, wordToInt[word])
			}

		}

		return wordInts
	}

	procFunc.WriteWordIntMappings = func(outputDir string) {

		wordToIntFn := filepath.Join(outputDir, wordToIntFilename)
		wordToIntF, err := os.OpenFile(wordToIntFn, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", wordToIntFn, err)
			os.Exit(-1)
		}

		intToWordFn := filepath.Join(outputDir, intToWordFilename)
		intToWordF, err := os.OpenFile(intToWordFn, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", intToWordFn, err)
			os.Exit(-1)
		}

		wordToIntBytes, err := json.Marshal(wordToInt)
		if err != nil {
			fmt.Printf("Error marshalling word to int\n")
			os.Exit(-1)
		}

		intToWordBytes, err := json.Marshal(intToWord)
		if err != nil {
			fmt.Printf("Error marshalling word to int\n")
			os.Exit(-1)
		}

		if _, err := wordToIntF.Write(wordToIntBytes); err != nil {
			fmt.Printf("Error writing to file %s: %v\n", wordToIntFn, err)
			os.Exit(-1)
		}

		if _, err := intToWordF.Write(intToWordBytes); err != nil {
			fmt.Printf("Error writing to file %s: %v\n", intToWordFn, err)
			os.Exit(-1)
		}
	}

	return &procFunc
}

func WordsToInts(stopWords []string, dataFilename string,
	outputDir string, wordIntsFn string) {

	proc := GenProcFunc(stopWords)

	wordIntsFilename := filepath.Join(outputDir, wordIntsFn)
	f, err := os.OpenFile(wordIntsFilename, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Printf("Error opening file %s, err: %v\n", wordIntsFilename, err)
		os.Exit(-1)
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	if bw == nil {
		fmt.Printf("Error creating new buffered writer\n")
		os.Exit(-1)
	}
	defer bw.Flush()

	in, done := data.LoadData(dataFilename)
	words := make([]string, maxWordsPerDoc)
	wordInts := make([]int, maxWordsPerDoc)
	var line string
LOOP:
	for {
		select {
		case line = <-in:
			line = proc.Replace(line)
			line = proc.ToLower(line)
			words = proc.GetWords(line, words)
			words = proc.RemoveStopwords(words)
			wordInts = proc.WordsToInts(words, wordInts)
			fmt.Println(wordInts)
			bw.WriteString(fmt.Sprintf("%v\n", wordInts))
		case <-done:
			break LOOP
		}
	}

	proc.WriteWordIntMappings(outputDir)
}
