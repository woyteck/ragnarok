package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Connect(host, user, pass, name string) *sql.DB {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v/%v?sslmode=disable", user, pass, host, name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
