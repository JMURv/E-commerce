package kafka

import (
	"github.com/IBM/sarama"
	conf "github.com/JMURv/e-commerce/favorites/pkg/config"
	"log"
)

type Broker struct {
	topic    string
	producer sarama.AsyncProducer
}

func New(conf *conf.KafkaConfig) *Broker {
	producer, err := sarama.NewAsyncProducer(conf.Addrs, nil)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	return &Broker{
		topic:    conf.NotificationTopic,
		producer: producer,
	}
}

func (b *Broker) Close() {
	if err := b.producer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}
}
