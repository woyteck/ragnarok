package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/scraper"
	"woyteck.pl/ragnarok/types"
)

const kafkaTopicScrap = "scrap_jobs"
const kafkaTopicIndex = "index_jobs"

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

	c.SubscribeTopics([]string{kafkaTopicScrap}, nil)

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
	var data types.ScrapTaskEvent
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		fmt.Printf("json unserialization error: %s", err)
		return
	}

	exists, _, err := store.Memory.GetMemoryBySource(context.Background(), data.Url)
	if err != nil {
		panic(err)
	}
	if exists {
		fmt.Println("already scrapped", data.Url)
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

	for _, paragraph := range paragraphs {
		fragment := types.NewMemoryFragment(paragraph, "", false, false, memory.ID)
		store.MemoryFragment.InsertMemoryFragment(context.Background(), fragment)
		emitIndexTask(fragment.ID)
	}
}

func emitIndexTask(memoryFragmentId uuid.UUID) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	task := types.IndexMemoryFragmentEvent{
		MemoryFragmentID: memoryFragmentId,
	}

	b, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}

	topic := kafkaTopicIndex
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: b,
	}, nil)
	if err != nil {
		panic(err)
	}

	p.Flush(15 * 1000)
}
