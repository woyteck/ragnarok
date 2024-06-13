package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/types"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	dbConn := db.Connect(host, user, pass, name)

	ctx := context.Background()
	conversationStore := db.NewPostgresConversationStore(dbConn, "conversations")

	conversationStore.Truncate(ctx)

	c := types.Conversation{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}
	conversationStore.InsertConversation(ctx, &c)
}
