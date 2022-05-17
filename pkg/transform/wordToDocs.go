package transform

import (
	"fmt"

	"github.com/find-in-docs/upload/pkg/data"
)

const (
	blockSize = 2
)

func wordToDocuments(docToWords map[data.DocumentId][]data.WordInt) map[data.WordInt][]data.DocumentId {

	var ok bool
	wordToDocs := make(map[data.WordInt][]data.DocumentId)
	var docIds []data.DocumentId

	for docId, wordIds := range docToWords {
		for _, word := range wordIds {
			docIds, ok = wordToDocs[word]

			if !ok {
				wordToDocs[word] = make([]data.DocumentId, 0)
				docIds, _ = wordToDocs[word]
			}

			wordToDocs[word] = append(docIds, docId)
		}
	}

	return wordToDocs
}

func WordToDocs(inputDocs <-chan *data.Doc,
	storeWordToDocMappings func(string, map[data.WordInt][]data.DocumentId) error) error {

	var err error
	docToWords := make(map[data.DocumentId][]data.WordInt)

	for {
		select {
		case doc, ok := <-inputDocs:

			if !ok { // channel closed by writer
				return nil
			}

			if doc == nil {
				// It is an error to get a doc that is nil.
				// It means, the channel has been closed.
				// Treat it like this by returning from this function.
				return nil
			}

			if len(docToWords) == blockSize {

				if err = storeWordToDocMappings(`wordid_to_docids`,
					wordToDocuments(docToWords)); err != nil {

					fmt.Printf("Error storing word-to-docs mappings: %v\n", err)
					return err
				} else {
					docToWords = make(map[data.DocumentId][]data.WordInt)
				}

				continue
			}

			if docToWords[doc.DocId] == nil {
				docToWords[doc.DocId] = make([]data.WordInt, len(doc.WordInts))
			}

			docToWords[doc.DocId] = doc.WordInts
		}
	}

	return nil
}
