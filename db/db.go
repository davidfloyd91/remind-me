package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	host   = "localhost"
	port   = 5432
	dbname = "remind_me"
)

func Connect() {
	info := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable", host, port, dbname)

	var err error
	DB, err = sql.Open("postgres", info)
	if err != nil {
		panic(err)
	}

	// defer DB.Close()

	err = DB.Ping()
	if err != nil {
		panic(err)
	}
}
