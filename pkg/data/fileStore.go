package data

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func StoreDataOnDisk(outputDir string, wordIntsFn string) (func([]int), func()) {

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

	return func(wordInts []int) {
			fmt.Fprintf(bw, "%v\n", wordInts)
		}, func() {
			bw.Flush()
			f.Close()
		}
}
