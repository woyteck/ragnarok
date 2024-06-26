package kafka

import (
	kaf "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gofiber/fiber/v2/log"
)

type Kafka struct {
	Servers string
}

func NewKafka(servers string) *Kafka {
	return &Kafka{
		Servers: servers,
	}
}

func (k *Kafka) createProducer() (*kaf.Producer, error) {
	producer, err := kaf.NewProducer(&kaf.ConfigMap{"bootstrap.servers": k.Servers})
	if err != nil {
		return nil, err
	}

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kaf.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Infof("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return producer, nil
}

func (k *Kafka) Produce(topic string, message []byte) error {
	producer, err := k.createProducer()
	if err != nil {
		return err
	}

	defer producer.Close()

	err = producer.Produce(&kaf.Message{
		TopicPartition: kaf.TopicPartition{
			Topic:     &topic,
			Partition: kaf.PartitionAny,
		},
		Value: message,
	}, nil)
	if err != nil {
		return err
	}

	producer.Flush(15 * 1000)

	return nil
}
