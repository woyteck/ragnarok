package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/joho/godotenv"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/openai"
	"woyteck.pl/ragnarok/prompter"
	"woyteck.pl/ragnarok/types"
	"woyteck.pl/ragnarok/vectordb"
)

const (
	kafkaTopicIndex   = "index_jobs"
	embeddingModel    = "text-embedding-ada-002"
	qdrantCollection  = "memory"
	embeddingSize     = 1536
	embeddingDistance = "Cosine"
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
	llm, ok := container.Get("llm").(*openai.Client)
	if !ok {
		panic("get llm failed")
	}
	prompter, ok := container.Get("prompter").(*prompter.Prompter)
	if !ok {
		panic("get prompter failed")
	}
	qdrant, ok := container.Get("vectordb").(*vectordb.QdrantClient)
	if !ok {
		panic("get vectordb failed")
	}

	err = createMemoryCollection(qdrant)
	if err != nil {
		panic(err)
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{kafkaTopicIndex}, nil)

	run := true

	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			err = indexMemoryFragment(msg, store, llm, prompter, qdrant)
			if err != nil {
				fmt.Printf("Indexer error: %v (%v)\n", err, msg)
			}
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}
	}

	c.Close()
}

func createMemoryCollection(qdrant *vectordb.QdrantClient) error {
	collectionExists, err := qdrant.CollectionExists(qdrantCollection)
	if err != nil {
		return err
	}

	if !collectionExists {
		err = qdrant.CreateCollection(qdrantCollection, embeddingSize, embeddingDistance)
		if err != nil {
			return err
		}
	}

	return nil
}

func indexMemoryFragment(msg *kafka.Message, store *db.Store, llm *openai.Client, prompter *prompter.Prompter, qdrant *vectordb.QdrantClient) error {
	var data types.IndexMemoryFragmentEvent
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		return err
	}

	memoryFragment, err := store.MemoryFragment.GetMemoryFragmentByUUID(context.Background(), data.MemoryFragmentID)
	if err != nil {
		return err
	}

	if !memoryFragment.IsRefined {
		contentRefined, err := generateRefinedContents(memoryFragment, llm, prompter)
		if err != nil {
			return err
		}

		if contentRefined != "nieistotne" && contentRefined != "" {
			memoryFragment.ContentRefined = contentRefined
		}
		memoryFragment.IsRefined = true

		err = store.MemoryFragment.UpdateMemoryFragment(context.Background(), memoryFragment)
		if err != nil {
			return err
		}
	}

	if !memoryFragment.IsEmbedded && memoryFragment.ContentRefined != "" {
		err = embedMemoryFragment(memoryFragment, llm, qdrant)
		if err != nil {
			return err
		}

		memoryFragment.IsEmbedded = true
		err = store.MemoryFragment.UpdateMemoryFragment(context.Background(), memoryFragment)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateRefinedContents(memoryFragment *types.MemoryFragment, llm *openai.Client, prompter *prompter.Prompter) (string, error) {
	llmContext, err := prompter.Get("summarize")
	if err != nil {
		return "", err
	}

	messages := []*openai.Message{
		{
			Role:    "system",
			Content: llmContext,
		},
		{
			Role:    "user",
			Content: memoryFragment.ContentOriginal,
		},
	}

	contentRefined, err := llm.GetCompletionShort(messages, "gpt-4-turbo")
	if err != nil {
		return "", err
	}

	return contentRefined, nil
}

func embedMemoryFragment(memoryFragment *types.MemoryFragment, llm *openai.Client, qdrant *vectordb.QdrantClient) error {
	vector, err := llm.GetEmbedding(memoryFragment.ContentRefined, "text-embedding-ada-002")
	if err != nil {
		return err
	}

	payload := map[string]any{
		"content":  memoryFragment.ContentRefined,
		"memoryID": memoryFragment.MemoryID,
	}

	return qdrant.UpsertPoints(qdrantCollection, vector, memoryFragment.ID, payload)
}
