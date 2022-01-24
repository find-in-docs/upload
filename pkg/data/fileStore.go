package data

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func StoreDataOnDisk(outputDir string, wordIntsFn string) (func(*Doc, []WordInt), func()) {

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

	return func(doc *Doc, wordInts []WordInt) {
			result := make([]WordInt, len(wordInts)+1)
			doc.WordInts = result
			copy(result[1:], wordInts)
			result[0] = doc.DocId
			fmt.Fprintf(bw, "%v\n", result)
		}, func() {
			bw.Flush()
			f.Close()
		}
}
