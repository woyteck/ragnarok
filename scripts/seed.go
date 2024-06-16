package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/types"
	"woyteck.pl/ragnarok/vectordb"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	container := di.NewContainer(di.Services)

	store, ok := container.Get("store").(*db.Store)
	if !ok {
		panic("get store failed")
	}
	qdrant, ok := container.Get("vectordb").(*vectordb.QdrantClient)
	if !ok {
		panic("get qdrant failed")
	}

	ctx := context.Background()

	store.Conversation.Truncate(ctx)
	store.Message.Truncate(ctx)

	now := time.Now()
	c := types.Conversation{
		ID:        uuid.New(),
		CreatedAt: &now,
	}
	store.Conversation.InsertConversation(ctx, &c)

	m := types.Message{
		ID:             uuid.New(),
		ConversationId: c.ID,
		Role:           "system",
		Content:        "Lorem ipsum",
		CreatedAt:      &now,
	}
	store.Message.InsertMessage(ctx, &m)

	store.MemoryFragment.Truncate(ctx)
	store.Memory.Truncate(ctx)

	qdrant.DeleteCollection("memory")
}
