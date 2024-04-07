package kafka

import (
	"github.com/IBM/sarama"
	_ "github.com/JMURv/e-commerce/query/pkg/config"
	"log"
)

type Broker struct {
	consumer sarama.Consumer
	cfg      *config.Config
}

func New(addr string) *Broker {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	brokers := []string{addr}

	consumer, err := sarama.NewConsumer(brokers, consumerConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	// Subscribe to Kafka topics
	topics := []string{"topic1", "topic2", "topic3", "topic4", "topic5", "topic6"}
	partitionConsumers := make(map[string]sarama.PartitionConsumer)
	for _, topic := range topics {
		partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error consuming Kafka topic %s: %v", topic, err)
		}
		partitionConsumers[topic] = partitionConsumer

		go func(topic string, partitionConsumer sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				switch msg.Topic {
				case s.cfg.KafkaTopics.ProductCreate.TopicName:
					s.processCreateProduct(ctx, r, msg)
				case s.cfg.KafkaTopics.ProductUpdate.TopicName:
					s.processUpdateProduct(ctx, r, msg)
				case s.cfg.KafkaTopics.ProductDelete.TopicName:
					s.processDeleteProduct(ctx, r, msg)
				}
				log.Printf("Received message from topic %s: %s", msg.Topic, string(msg.Value))
			}
		}(topic, partitionConsumer)
	}
	return &Broker{consumer: consumer}
}

func (b *Broker) Start() {

}

func (b *Broker) Close() {
	if err := b.consumer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}
}
