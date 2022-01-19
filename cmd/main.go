package main

import (
	"fmt"
	"path/filepath"

	"github.com/samirgadkari/search/pkg/config"
	"github.com/samirgadkari/search/pkg/transform"
)

type Config struct {
	RawDocumentsFn string `yaml:"raw-documents-fn"`
	OutputDirFn    string `yaml:"output-dir-fn"`
	StopwordsFn    string `yaml:"stopwords-fn"`
}

type Stopwords struct {
	English []string `json:"english_stopwords"`
}

func main() {

	cfg := config.LoadConfig()
	fmt.Printf("%#v\n", cfg)

	stopWords := config.LoadStopwords(cfg)
	fmt.Println(stopWords)

	transform.WordsToInts(stopWords, cfg.DataFile,
		filepath.Join(cfg.OutputDir, cfg.WordIntsFile))
}
