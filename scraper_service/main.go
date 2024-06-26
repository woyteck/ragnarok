package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	kaf "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/joho/godotenv"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/indexer"
	"woyteck.pl/ragnarok/kafka"
	"woyteck.pl/ragnarok/scraper"
	"woyteck.pl/ragnarok/types"
)

const kafkaTopicScrap = "scrap_jobs"
const kafkaTopicIndex = "index_jobs"

type ScraperConsumer struct {
	scraper *scraper.CollyScraper
	store   *db.Store
	indexer *indexer.Indexer
	kafka   *kafka.Kafka
}

func NewScraperConsumer(container *di.Container, kafka *kafka.Kafka) *ScraperConsumer {
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

	return &ScraperConsumer{
		scraper: scraper,
		store:   store,
		indexer: indexer,
		kafka:   kafka,
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	container := di.NewContainer(di.Services)

	kafka, ok := container.Get("kafka").(*kafka.Kafka)
	if !ok {
		panic("get kafka failed")
	}

	topic := kafkaTopicScrap
	consumer := NewScraperConsumer(container, kafka)
	kafka.Consume(topic, consumer.saveMemory)
}

func (c *ScraperConsumer) saveMemory(msg *kaf.Message) {
	var data types.ScrapTaskEvent
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		fmt.Printf("json unserialization error: %s", err)
		return
	}

	exists, memory, err := c.store.Memory.GetMemoryBySource(context.Background(), data.Url)
	if err != nil {
		panic(err)
	}
	if exists && memory.Content != "" {
		fmt.Println("already scrapped", data.Url)
		return
	}

	title, html, err := c.scraper.GetArticle(data.Url, data.CssSelector)
	if err != nil {
		fmt.Printf("failed to scrap url: %s, css selector: %s, error: %s", data.Url, data.CssSelector, err)
		return
	}

	if exists {
		memory.Content = html
	} else {
		memory = types.NewMemory(types.MemoryTypeWebArticle, data.Url, html)
	}

	insertedMemoryFragmentIds, err := c.indexer.Index(memory, !exists, title, data.Url)
	if err != nil {
		fmt.Printf("indexer error: %s", err)
		return
	}

	for _, id := range insertedMemoryFragmentIds {
		topic := kafkaTopicIndex
		task := types.IndexMemoryFragmentEvent{
			MemoryFragmentID: id,
		}

		message, err := json.Marshal(task)
		if err != nil {
			panic(err)
		}
		c.kafka.Produce(topic, message)
	}
}
