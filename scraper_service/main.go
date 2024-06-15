package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gofiber/fiber/v2/log"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/scraper"
)

const kafkaTopic = "scrap_jobs"

type ScrapTask struct {
	Url         string `json:"url"`
	CssSelector string `json:"cssSelector"`
}

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	fmt.Println(os.Args[1])

	container := di.NewContainer(di.Services)

	scraper, ok := container.Get("scraper").(*scraper.CollyScraper)
	if !ok {
		panic("get scraper failed")
	}

	_ = scraper

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{kafkaTopic}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			log.Errorf("kafka consume error: %s", err)
			continue
		}

		var data ScrapTask
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			log.Errorf("json unserialization error: %s", err)
			continue
		}

		text, err := scraper.ScrapPage(data.Url, data.CssSelector)
		if err != nil {
			log.Errorf("failed to scrap url: %s, css selector: %s, error: %s", data.Url, data.CssSelector, err)
			continue
		}

		fmt.Println(text) //TODO: store text to db, emit event for indexer service (also TODO)
	}
}
