package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type Truncater interface {
	Truncate(context.Context) error
}

func Connect(host, user, pass, name string) *sql.DB {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v/%v?sslmode=disable", user, pass, host, name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func makePlaceholders(count int) string {
	list := []string{}
	for i := 0; i < count; i++ {
		list = append(list, fmt.Sprintf("$%d", i+1))
	}

	return strings.Join(list, ",")
}
