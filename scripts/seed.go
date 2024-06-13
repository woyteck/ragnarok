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
	messageStore := db.NewPostgresMessagesStore(dbConn, "messages")

	messageStore.Truncate(ctx)
	conversationStore.Truncate(ctx)

	now := time.Now()
	c := types.Conversation{
		ID:        uuid.New(),
		CreatedAt: &now,
	}
	conversationStore.InsertConversation(ctx, &c)

	m := types.Message{
		ID:             uuid.New(),
		ConversationId: c.ID,
		Role:           "system",
		Content:        "Lorem ipsum",
		CreatedAt:      &now,
	}
	messageStore.InsertMessage(ctx, &m)
}
