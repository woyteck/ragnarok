package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"woyteck.pl/ragnarok/types"
)

const kafkaTopicScrap = "scrap_jobs"

func main() {
	url := flag.String("url", "", "url to scrap")
	cssSelector := flag.String("selector", ".article-content p", "css selector to extract paragraphs of text")
	flag.Parse()

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
