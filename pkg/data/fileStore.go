package data

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func StoreDataOnDisk(outputDir string, wordIntsFn string, in <-chan []int) {

	go func() {
		wordIntsFilename := filepath.Join(outputDir, wordIntsFn)
		f, err := os.Create(wordIntsFilename)
		if err != nil {
			fmt.Printf("Error opening file %s, err: %v\n", wordIntsFilename, err)
			os.Exit(-1)
		}
		defer f.Close() // This defer is first, so it will run last

		bw := bufio.NewWriterSize(f, 256)
		if bw == nil {
			fmt.Printf("Error creating new buffered writer\n")
			os.Exit(-1)
		}
		defer bw.Flush() // This defer is second, so it will run first

		for wordInts := range in {

			fmt.Fprintf(bw, "%v\n", wordInts)
		}
	}()
}
