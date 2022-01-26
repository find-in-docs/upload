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
	OpenConnection  func() *ent.Client
	CreateSchema    func(*ent.Client)
	StoreData       func(*ent.Client, *Doc)
	ReadData        func(*ent.Client) *Doc
	CloseConnection func(*ent.Client)
}

func DBSetup() *DBFunc {

	var dbFunc DBFunc

	dbFunc.OpenConnection = func() *ent.Client {

		client, err := ent.Open("postgres", viper.GetString("output.location"))
		if err != nil {
			fmt.Printf("Error connecting to postgres database using: %s\n",
				viper.GetString("output.location"))
			fmt.Printf("err: %v\n", err)
			os.Exit(-1)
		}

		return client
	}
	dbFunc.CreateSchema = func(client *ent.Client) {

		if client == nil {
			fmt.Printf("Create client connection before creating schema")
			os.Exit(-1)
		}

		if err := client.Schema.Create(context.Background()); err != nil {
			fmt.Printf("Failed to create the schema. err: %v", err)
			os.Exit(-1)
		}
	}
	dbFunc.StoreData = func(client *ent.Client, doc *Doc) {}
	dbFunc.ReadData = func(client *ent.Client) *Doc { return nil }
	dbFunc.CloseConnection = func(client *ent.Client) {

		if client == nil {
			fmt.Printf("client is nil\n")
			os.Exit(-1)
		}

		client.Close()
	}

	return &dbFunc
}
