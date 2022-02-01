package data

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func ExecuteQueryRows(tableName string, query string, params ...interface{}) (pgx.Rows, error) {

	var err error
	var db *DB

	db, err = DBConnect()
	if err != nil {
		return nil, err
	}

	var rows pgx.Rows
	rows, err = db.conn.Query(context.Background(),
		"select * from "+tableName+";")

	return rows, err
}

func LoadWordToIntTable(tableName string) (*map[string]WordInt, error) {

	var err error
	var rows pgx.Rows

	rows, err = ExecuteQueryRows("wordtoint", "select * from wordtoint;")
	if err != nil {
		return nil, err
	}

	WordToInt := make(map[string]WordInt)

	var word string
	var int_ WordInt

	for rows.Next() {
		err = rows.Scan(&word, &int_)
		if err != nil {
			fmt.Printf("Error loading wordtoint table: %v\n", err)
			return nil, err
		}

		WordToInt[word] = int_
	}

	// This catches errors reported by rows.Next()
	if rows.Err() != nil {
		return nil, err
	}
	return &WordToInt, nil
}

func LoadDocs() (<-chan *Doc, error) {

	var err error
	var rows pgx.Rows

	rows, err = ExecuteQueryRows("doc", "select * from doc;")
	if err != nil {
		return nil, err
	}

	inputDocs := make(chan *Doc)

	go func() {
		for rows.Next() {

			var docId uint64
			var wordInts []uint64

			doc := new(Doc)
			err = rows.Scan(&docId, &wordInts, nil, nil, nil, nil,
				nil, nil, nil, nil, nil)
			if err != nil {
				fmt.Printf("Error loading doc table: %v\n", err)
				return
			}

			doc.DocId = DocumentId(docId)
			doc.WordInts = make([]WordInt, len(wordInts))
			for i, v := range wordInts {
				doc.WordInts[i] = WordInt(v)
			}

			inputDocs <- doc
		}

		// This catches errors reported by rows.Next()
		if rows.Err() != nil {
			fmt.Printf("rows.Next() reported an error: %v\n", err)
			return
		}

		close(inputDocs)
	}()

	return inputDocs, nil
}
