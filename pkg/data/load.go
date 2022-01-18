package data

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func LoadData(dataFile *string) (<-chan string, <-chan struct{}) {
	done := make(chan struct{})
	in := make(chan string, 100)

	go func() {
		fmt.Println("In LoadData")
		f, err := os.Open(*dataFile)
		defer f.Close()

		if err != nil {
			log.Fatalf("Error opening file %s", *dataFile)
			os.Exit(-1)
		}

		s := bufio.NewScanner(f)
		s.Split(bufio.ScanLines)

		for s.Scan() {
			in <- s.Text()
		}

		close(in)
		close(done)
	}()

	return in, done
}
