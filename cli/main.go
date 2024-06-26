package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/joho/godotenv"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/indexer"
	"woyteck.pl/ragnarok/scraper"
	"woyteck.pl/ragnarok/types"
)

const kafkaTopicScrap = "scrap_jobs"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := flag.String("url", "", "url to scrap")
	cssSelector := flag.String("selector", ".article-content p", "css selector to extract paragraphs of text")
	mode := flag.String("mode", "", "queue | experiment")
	flag.Parse()

	if *mode == "experiment" {
		experiment()
	} else {
		if *url == "" {
			fmt.Fprintf(os.Stderr, "url is required\n")
			return

		}
		if *cssSelector == "" {
			fmt.Fprintf(os.Stderr, "cssSelector is required\n")
			return
		}

		emitScrapTask(*url, *cssSelector)
	}
}

func experiment() {
	container := di.NewContainer(di.Services)

	url := "https://italia-by-natalia.pl/etna-jak-zwiedzac-wulkan-informacje-praktyczne/"
	selector := ".article-content"

	scraper, ok := container.Get("scraper").(*scraper.CollyScraper)
	if !ok {
		panic("get scraper failed")
	}

	title, html, err := scraper.GetArticle(url, selector)
	if err != nil {
		fmt.Printf("failed to scrap url: %s, css selector: %s, error: %s", url, selector, err)
		return
	}

	indexer, ok := container.Get("indexer").(*indexer.Indexer)
	if !ok {
		panic("get indexer failed")
	}

	err = indexer.Index(html, title)
	if err != nil {
		panic(err)
	}
}

func emitScrapTask(url string, cssSelector string) {
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

	task := types.ScrapTaskEvent{
		Url:         url,
		CssSelector: cssSelector,
	}

	b, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}

	topic := kafkaTopicScrap
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
