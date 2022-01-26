package data

import (
	"context"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/samirgadkari/search/ent"
	"github.com/spf13/viper"
)

type DBFunc struct {
	OpenConnection  func()
	CreateSchema    func()
	StoreData       func(*Doc)
	ReadData        func() *Doc
	CloseConnection func()
}

func DBSetup() *DBFunc {

	var client *ent.Client
	var err error
	var dbFunc DBFunc

	dbFunc.OpenConnection = func() {

		client, err = ent.Open("postgres", viper.GetString("output.location"))
		if err != nil {
			fmt.Printf("Error connecting to postgres database using: %s\n",
				viper.GetString("output.location"))
			fmt.Printf("err: %v\n", err)
			os.Exit(-1)
		}
	}
	dbFunc.CreateSchema = func() {

		if client == nil {
			fmt.Printf("Create client connection before creating schema")
			os.Exit(-1)
		}

		if err := client.Schema.Create(context.Background()); err != nil {
			fmt.Printf("Failed to create the schema. err: %v", err)
			os.Exit(-1)
		}
	}
	dbFunc.StoreData = func(doc *Doc) {}
	dbFunc.ReadData = func() *Doc { return nil }
	dbFunc.CloseConnection = func() {

		if client == nil {
			fmt.Printf("client is nil\n")
			os.Exit(-1)
		}

		client.Close()
	}

	return &dbFunc
}
