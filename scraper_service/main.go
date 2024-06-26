package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	kaf "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/indexer"
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

	indexer, ok := container.Get("indexer").(*indexer.Indexer)
	if !ok {
		panic("get indexer failed")
	}

	c, err := kaf.NewConsumer(&kaf.ConfigMap{
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
			saveMemory(msg, scraper, store, indexer)
		} else if !err.(kaf.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}
	}

	c.Close()
}

func saveMemory(msg *kaf.Message, scraper *scraper.CollyScraper, store *db.Store, indexer *indexer.Indexer) {
	var data types.ScrapTaskEvent
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		fmt.Printf("json unserialization error: %s", err)
		return
	}

	exists, memory, err := store.Memory.GetMemoryBySource(context.Background(), data.Url)
	if err != nil {
		panic(err)
	}
	if exists && memory.Content != "" {
		fmt.Println("already scrapped", data.Url)
		return
	}

	title, html, err := scraper.GetArticle(data.Url, data.CssSelector)
	if err != nil {
		fmt.Printf("failed to scrap url: %s, css selector: %s, error: %s", data.Url, data.CssSelector, err)
		return
	}

	if exists {
		memory.Content = html
	} else {
		memory = types.NewMemory(types.MemoryTypeWebArticle, data.Url, html)
	}

	insertedMemoryFragmentIds, err := indexer.Index(memory, !exists, title, data.Url)
	if err != nil {
		fmt.Printf("indexer error: %s", err)
		return
	}

	for _, id := range insertedMemoryFragmentIds {
		emitIndexTask(id)
	}
}

func emitIndexTask(memoryFragmentId uuid.UUID) {
	p, err := kaf.NewProducer(&kaf.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kaf.Message:
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
	err = p.Produce(&kaf.Message{
		TopicPartition: kaf.TopicPartition{
			Topic:     &topic,
			Partition: kaf.PartitionAny,
		},
		Value: b,
	}, nil)
	if err != nil {
		panic(err)
	}

	p.Flush(15 * 1000)
}
