package kafka

import (
	"github.com/IBM/sarama"
	conf "github.com/JMURv/e-commerce/notifications/pkg/config"
	"log"
)

type Broker struct {
	cfg      *conf.Config
	consumer sarama.Consumer
}

// Start Kafka consumer loop
//go func() {
//	for {
//		m, err := reader.ReadMessage(ctx)
//		if err != nil && errors.Is(err, io.EOF) {
//			log.Printf("Kafka has been stopped")
//			return
//		} else if err != nil {
//			log.Printf("Error reading message from Kafka: %v", err)
//			continue
//		}
//		log.Println(m.Topic)
//		log.Printf("Received message from Kafka: %s", m.Value)
//	}
//}()

func New(conf *conf.Config) *Broker {
	consumer, err := sarama.NewConsumer(conf.Kafka.Addrs, nil)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	return &Broker{
		cfg:      conf,
		consumer: consumer,
	}
}

func (b *Broker) Close() {
	if err := b.consumer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}
}
