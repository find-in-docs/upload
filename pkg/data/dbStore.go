package data

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
)

type DBFunc struct {
	OpenConnection       func()
	CreateTable          func(string)
	StoreData            func(*Doc, string, []WordInt)
	StoreWordIntMappings func(string, map[string]WordInt, string, map[WordInt]string)
	ReadData             func() *Doc
	CloseConnection      func()
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

	var conn *pgx.Conn
	var err error
	var dbFunc DBFunc
	docSchema := `(docid integer,
			wordints integer[],
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
				int integer)`
	createWordToIntString := `(word,
					int)`
	intToWordSchema := `(int integer, 
				word text)`
	createIntToWordString := `(int,
					word)`

	dbFunc.OpenConnection = func() {

		conn, err = pgx.Connect(context.Background(), viper.GetString("output.connection"))
		if err != nil {
			fmt.Printf("Error connecting to postgres database using: %s\n",
				viper.GetString("output.location"))
			fmt.Printf("err: %v\n", err)
			os.Exit(-1)
		}
	}
	dbFunc.CreateTable = func(tableName string) {

		if conn == nil {
			fmt.Printf("Create db connection before creating schema")
			os.Exit(-1)
		}

		createTable(conn, tableName, docSchema)
	}
	dbFunc.StoreWordIntMappings = func(wordToIntTable string, wordToInt map[string]WordInt,
		intToWordTable string, intToWord map[WordInt]string) {

		createTable(conn, wordToIntTable, wordToIntSchema)
		createTable(conn, intToWordTable, intToWordSchema)

		wordToIntInsertStatement := `insert into ` + wordToIntTable + ` ` + createWordToIntString +
			`values ($1, $2);`
		intToWordInsertStatement := `insert into ` + intToWordTable + ` ` + createIntToWordString +
			`values ($1, $2);`
		for word, i := range wordToInt {
			if _, err := conn.Exec(context.Background(), wordToIntInsertStatement,
				word, i); err != nil {
				fmt.Printf("Store Int to Word mapping failed. err: %v\n", err)
				os.Exit(-1)
			}
			if _, err := conn.Exec(context.Background(), intToWordInsertStatement,
				i, word); err != nil {
				fmt.Printf("Store Word to Int mapping failed. err: %v\n", err)
				os.Exit(-1)
			}
		}

	}
	dbFunc.StoreData = func(doc *Doc, tableName string, wordInts []WordInt) {

		insertStatement := `insert into ` + tableName + ` ` + createDocString +
			` values ($1, $2, $3, $4, $5, 
			 $6, $7, $8, $9, $10, $11);`

		if _, err := conn.Exec(context.Background(), insertStatement,
			doc.DocId, doc.WordInts, doc.InputDocId,
			doc.UserId, doc.BusinessId, doc.Stars, doc.Useful,
			doc.Funny, doc.Cool, doc.Text, doc.Date); err != nil {
			fmt.Printf("Store data failed. err: %v\n", err)
			os.Exit(-1)
		}
	}
	dbFunc.ReadData = func() *Doc { return nil }
	dbFunc.CloseConnection = func() {

		if conn == nil {
			fmt.Printf("conn is nil\n")
			os.Exit(-1)
		}

		err := conn.Close(context.Background())
		if err != nil {
			fmt.Printf("Error closing DB connection: %v\n", err)
			os.Exit(-1)
		}
	}

	return &dbFunc
}
