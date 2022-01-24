package data

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type WordInt uint32

func StoreDataOnDisk(outputDir string, wordIntsFn string) (func([]WordInt), func()) {

	var docId WordInt = 0

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

	return func(wordInts []WordInt) {
			result := make([]WordInt, len(wordInts)+1)
			copy(result[1:], wordInts)
			result[0] = docId
			fmt.Fprintf(bw, "%v\n", result)
			docId += 1
		}, func() {
			bw.Flush()
			f.Close()
		}
}
