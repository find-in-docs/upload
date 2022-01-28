package data

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
)

type DB struct {
	conn *pgx.Conn
}

func DBConnect() (*DB, error) {

	conn, err := pgx.Connect(context.Background(), viper.GetString("output.connection"))
	if err != nil {
		fmt.Printf("Error connecting to postgres database using: %s\n",
			viper.GetString("output.location"))
		fmt.Printf("err: %v\n", err)
		return nil, err
	}

	db := DB{conn}

	return &db, nil
}

func (db *DB) CreateTable(tableName string, schema string) error {

	if db.conn == nil {
		fmt.Printf("Create db connection before creating schema\n")
		return fmt.Errorf("Create db connection before creating schema\n")
	}

	checkIfExists := `select 'public.` + tableName + `'::regclass;`
	if _, err := db.conn.Exec(context.Background(), checkIfExists); err != nil {
		fmt.Printf("Table %s does not exist, so create it.\n", tableName)

		createString := `create table ` + tableName + ` ` + schema + `;`
		if _, err := db.conn.Exec(context.Background(), createString); err != nil {
			fmt.Printf("Failed to create the schema. err: %v\n", err)
			return err
		}
	}

	return nil
}

func (db *DB) DBDisconnect() error {

	if db.conn == nil {
		fmt.Printf("conn is nil\n")
		os.Exit(-1)
	}

	err := db.conn.Close(context.Background())
	if err != nil {
		fmt.Printf("Error closing DB connection: %v\n", err)
		return err
	}

	return nil
}
