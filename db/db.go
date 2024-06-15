package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

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

func parseTimestamp(value sql.NullString) (*time.Time, error) {
	const timeFormat = "2006-01-02T15:04:05Z"
	createdAtTime, err := time.Parse(timeFormat, value.String)
	return &createdAtTime, err
}

type Store struct {
	Conversation ConversationStore
	Message      MessagesStore
	Cache        CacheStore
	Memory       MemoriesStore
}
