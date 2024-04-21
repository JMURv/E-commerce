package kafka

import (
	"github.com/IBM/sarama"
	cfg "github.com/JMURv/e-commerce/users/pkg/config"
	"log"
)

type Broker struct {
	producer sarama.AsyncProducer
	topic    string
}

func New(conf *cfg.KafkaConfig) *Broker {
	producer, err := sarama.NewAsyncProducer(conf.Addrs, nil)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	return &Broker{
		producer: producer,
		topic:    conf.NotificationTopic,
	}
}

func (b *Broker) Close() {
	if err := b.producer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}
}
