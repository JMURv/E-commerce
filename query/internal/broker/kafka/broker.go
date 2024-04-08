package kafka

import (
	"context"
	"github.com/IBM/sarama"
	conf "github.com/JMURv/e-commerce/query/pkg/config"
	"log"
)

type Broker struct {
	cfg                *conf.Config
	consumer           sarama.Consumer
	partitionConsumers map[string]sarama.PartitionConsumer
}

func New(addrs []string, conf *conf.Config) *Broker {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(addrs, consumerConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	return &Broker{
		cfg:                conf,
		consumer:           consumer,
		partitionConsumers: make(map[string]sarama.PartitionConsumer),
	}
}

func (b *Broker) Start() {
	ctx := context.Background()
	for _, topic := range b.cfg.Kafka.Topics {
		partitionConsumer, err := b.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error consuming Kafka topic %s: %v", topic, err)
		}
		b.partitionConsumers[topic] = partitionConsumer

		go func(topic string, partitionConsumer sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				switch msg.Topic {
				case b.cfg.Kafka.TopicName.Create:
					b.processCreate(ctx, topic, msg)
				case b.cfg.Kafka.TopicName.Update:
					b.processUpdate(ctx, topic, msg)
				case b.cfg.Kafka.TopicName.Delete:
					b.processDelete(ctx, topic, msg)
				}
				log.Printf("Received message from topic %s: %s", msg.Topic, string(msg.Value))
			}
		}(topic, partitionConsumer)
	}
}

func (b *Broker) Close() {
	// Close all partition consumers
	for _, topic := range b.cfg.Kafka.Topics {
		if err := b.partitionConsumers[topic].Close(); err != nil {
			log.Printf("Error closing Kafka partition consumer for topic %s: %v", topic, err)
		}
	}
	if err := b.consumer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}
}
