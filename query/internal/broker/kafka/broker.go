package kafka

import (
	"context"
	"github.com/IBM/sarama"
	conf "github.com/JMURv/e-commerce/query/pkg/config"
	"log"
)

type Broker struct {
	cfg      *conf.Config
	consumer sarama.Consumer
	pc       map[string]sarama.PartitionConsumer
}

func New(addrs []string, conf *conf.Config) *Broker {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(addrs, consumerConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	return &Broker{
		cfg:      conf,
		consumer: consumer,
		pc:       make(map[string]sarama.PartitionConsumer),
	}
}

func (b *Broker) Start() {
	ctx := context.Background()
	for _, topic := range b.cfg.Kafka.Topics {
		partitionConsumer, err := b.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error consuming Kafka topic %s: %v", topic, err)
		}
		b.pc[topic] = partitionConsumer

		go func(topic string, partitionConsumer sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				msgType := extractTypeField(msg.Headers)
				switch msgType {
				case b.cfg.Kafka.TopicName.Create:
					b.processCreate(ctx, topic, msgType, msg.Value)
				case b.cfg.Kafka.TopicName.Update:
					b.processUpdate(ctx, topic, msgType, msg.Value)
				case b.cfg.Kafka.TopicName.Delete:
					b.processDelete(ctx, topic, msgType, msg.Value)
				}
				log.Printf("Received message from topic %s: %s", msg.Topic, string(msg.Value))
			}
		}(topic, partitionConsumer)
	}
}

func extractTypeField(headers []*sarama.RecordHeader) string {
	for _, h := range headers {
		if string(h.Key) == "type" {
			return string(h.Value)
		}
	}
	return ""
}

func (b *Broker) Close() {
	// Close all partition consumers
	for _, topic := range b.cfg.Kafka.Topics {
		if err := b.pc[topic].Close(); err != nil {
			log.Printf("Error closing Kafka partition consumer for topic %s: %v", topic, err)
		}
	}
	if err := b.consumer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}
}

func (b *Broker) processCreate(ctx context.Context, topic string, msgType string, msg []byte) {
	switch topic {
	case "users":
		log.Println("users create request")
	}
}

func (b *Broker) processUpdate(ctx context.Context, topic string, msgType string, msg []byte) {}

func (b *Broker) processDelete(ctx context.Context, topic string, msgType string, msg []byte) {}
