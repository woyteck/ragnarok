package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/joho/godotenv"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/scraper"
	"woyteck.pl/ragnarok/types"
)

const kafkaTopic = "scrap_jobs"

type ScrapTask struct {
	Url         string `json:"url"`
	CssSelector string `json:"cssSelector"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	container := di.NewContainer(di.Services)

	scraper, ok := container.Get("scraper").(*scraper.CollyScraper)
	if !ok {
		panic("get scraper failed")
	}

	store, ok := container.Get("store").(*db.Store)
	if !ok {
		panic("get store failed")
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{kafkaTopic}, nil)

	run := true

	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			saveMemory(msg, scraper, store)
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

func saveMemory(msg *kafka.Message, scraper *scraper.CollyScraper, store *db.Store) {
	var data ScrapTask
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		fmt.Printf("json unserialization error: %s", err)
		return
	}

	paragraphs, err := scraper.ScrapPage(data.Url, data.CssSelector)
	if err != nil {
		fmt.Printf("failed to scrap url: %s, css selector: %s, error: %s", data.Url, data.CssSelector, err)
		return
	}

	text := strings.Join(paragraphs, "\n")
	memory := types.NewMemory(types.MemoryTypeWebArticle, data.Url, text)
	store.Memory.InsertMemory(context.Background(), memory)
}
