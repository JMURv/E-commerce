package kafka

import (
	"github.com/IBM/sarama"
	conf "github.com/JMURv/e-commerce/users/pkg/config"
	"log"
)

type Broker struct {
	cfg      *conf.Config
	producer sarama.AsyncProducer
}

func New(conf *conf.Config) *Broker {
	producer, err := sarama.NewAsyncProducer(conf.Kafka.Addrs, nil)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	return &Broker{
		cfg:      conf,
		producer: producer,
	}
}

func (b *Broker) Close() {
	if err := b.producer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}
}
