package data

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

type DBFunc struct {
	OpenConnection         func() error
	CreateTable            func(string) error
	StoreData              func(*Doc, string, []WordInt) error
	StoreWordIntMappings   func(string, map[string]WordInt) error
	StoreWordToDocMappings func(string, map[WordInt][]DocumentId) error
	ReadData               func() *Doc
	CloseConnection        func() error
}

func createTable(conn *pgx.Conn, tableName string, schema string) {

	checkIfExists := `select 'public.` + tableName + `'::regclass;`
	if _, err := conn.Exec(context.Background(), checkIfExists); err != nil {
		fmt.Printf("Table %s does not exist, so create it.\n", tableName)

		createString := `create table ` + tableName + ` ` + schema + `;`
		if _, err := conn.Exec(context.Background(), createString); err != nil {
			fmt.Printf("Failed to create the schema. err: %v\n", err)
			os.Exit(-1)
		}
	}
}

func DBSetup() *DBFunc {

	var db *DB
	var err error
	var dbFunc DBFunc
	docSchema := `(docid bigint,
			wordints bigint[],
			inputdocId varchar(25),
			userid varchar(25),
			businessId varchar(25),
			stars real, 
			useful smallint,
			funny smallint,
			cool smallint,
			text text,
			date varchar(25))`
	createDocString := `(docid,
			wordints,
			inputdocId,
			userid,
			businessId,
			stars, 
			useful,
			funny,
			cool,
			text,
			date)`
	wordToIntSchema := `(word text,
				int bigint)`
	createWordToIntString := `(word,
					int)`
	wordIdsToDocIdsSchema := `(wordid bigint, docids bigint[])`

	dbFunc.OpenConnection = func() error {

		db, err = DBConnect()
		if err != nil {
			return err
		}

		return nil
	}
	dbFunc.CreateTable = func(tableName string) error {

		return db.CreateTable(tableName, docSchema)
	}
	dbFunc.StoreWordIntMappings = func(wordToIntTable string, wordToInt map[string]WordInt) error {

		db.CreateTable(wordToIntTable, wordToIntSchema)

		wordToIntInsertStatement := `insert into ` + wordToIntTable + ` ` + createWordToIntString +
			`values ($1, $2);`
		for word, i := range wordToInt {
			if _, err := db.conn.Exec(context.Background(), wordToIntInsertStatement,
				word, i); err != nil {
				fmt.Printf("Store Int to Word mapping failed. err: %v\n", err)
				return err
			}
		}

		return nil
	}
	dbFunc.StoreData = func(doc *Doc, tableName string, wordInts []WordInt) error {

		insertStatement := `insert into ` + tableName + ` ` + createDocString +
			` values ($1, $2, $3, $4, $5, 
			 $6, $7, $8, $9, $10, $11);`

		if _, err := db.conn.Exec(context.Background(), insertStatement,
			doc.DocId, doc.WordInts, doc.InputDocId,
			doc.UserId, doc.BusinessId, doc.Stars, doc.Useful,
			doc.Funny, doc.Cool, doc.Text, doc.Date); err != nil {
			fmt.Printf("Store data failed. err: %v\n", err)
			return err
		}

		return nil
	}
	dbFunc.StoreWordToDocMappings = func(wordIdsToDocIdsTable string,
		wordToDocs map[WordInt][]DocumentId) error {

		db.CreateTable(wordIdsToDocIdsTable, wordIdsToDocIdsSchema)

		// In this update statement, the excluded docids are the ones that were not
		// inserted in.
		updateStatement := `
			insert into wordid_to_docids(wordid, docids) values($1, $2)
			on conflict(wordid) do
			update set docids=array(select distinct unnest(wordid_to_docids.docids || excluded.docids));
			`

		for k, v := range wordToDocs {

			if _, err := db.conn.Exec(context.Background(), updateStatement,
				k, v); err != nil {
				fmt.Printf("Update failed. err: %v\n", err)
				return err
			}
		}

		return nil
	}
	dbFunc.ReadData = func() *Doc { return nil }
	dbFunc.CloseConnection = func() error {

		return db.DBDisconnect()
	}

	return &dbFunc
}
