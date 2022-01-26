package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type DiskFunc struct {
	LoadStopwords        func() []string
	LoadDoc              func() (*Doc, bool)
	StoreData            func(*Doc, []WordInt)
	WriteWordIntMappings func(map[string]WordInt, map[WordInt]string)
	Close                func()
}

func DiskSetup(outputDir string, wordIntsFn string) *DiskFunc {

	wordToIntFilename := viper.GetString("output.wordToIntFn")
	intToWordFilename := viper.GetString("output.intToWordFn")
	dataFilename := viper.GetString("dataFile")

	wordIntsFilename := filepath.Join(outputDir, wordIntsFn)
	f, err := os.Create(wordIntsFilename)
	if err != nil {
		fmt.Printf("Error opening file %s, err: %v\n", wordIntsFilename, err)
		os.Exit(-1)
	}

	bw := bufio.NewWriter(f)
	if bw == nil {
		fmt.Printf("Error creating new buffered writer\n")
		os.Exit(-1)
	}

	var diskFunc DiskFunc
	diskFunc.LoadStopwords = func() []string {

		return Load(viper.GetString("englishStopwordsFile"))
	}

	diskFunc.LoadDoc = LoadDocFn(dataFilename)

	diskFunc.StoreData = func(doc *Doc, wordInts []WordInt) {
		result := make([]WordInt, len(wordInts)+1)
		doc.WordInts = result
		copy(result[1:], wordInts)
		result[0] = doc.DocId
		fmt.Fprintf(bw, "%v\n", result)
	}

	diskFunc.WriteWordIntMappings = func(wordToInt map[string]WordInt, intToWord map[WordInt]string) {

		wordToIntFn := filepath.Join(outputDir, wordToIntFilename)
		wordToIntF, err := os.Create(wordToIntFn)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", wordToIntFn, err)
			os.Exit(-1)
		}
		defer wordToIntF.Close()

		intToWordFn := filepath.Join(outputDir, intToWordFilename)
		intToWordF, err := os.Create(intToWordFn)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", intToWordFn, err)
			os.Exit(-1)
		}
		defer intToWordF.Close()

		wordToIntBytes, err := json.Marshal(wordToInt)
		if err != nil {
			fmt.Printf("Error marshalling word to int\n")
			os.Exit(-1)
		}

		intToWordBytes, err := json.Marshal(intToWord)
		if err != nil {
			fmt.Printf("Error marshalling int to word\n")
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
	diskFunc.Close = func() {
		bw.Flush()
		f.Close()
	}

	return &diskFunc
}
