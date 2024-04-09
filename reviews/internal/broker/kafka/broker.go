package kafka

import (
	"github.com/IBM/sarama"
	conf "github.com/JMURv/e-commerce/reviews/pkg/config"
	"log"
	"strconv"
)

type Broker struct {
	cfg      *conf.Config
	producer sarama.AsyncProducer
}

func New(addrs []string, conf *conf.Config) *Broker {
	producer, err := sarama.NewAsyncProducer(addrs, nil)
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

func (b *Broker) NewReviewNotification(reviewID uint64, msg []byte) error {
	b.producer.Input() <- &sarama.ProducerMessage{
		Topic: b.cfg.Kafka.NotificationTopic,
		Key:   sarama.StringEncoder(strconv.FormatUint(reviewID, 10)),
		Value: sarama.ByteEncoder(msg),
	}
	log.Println("Message sent successfully!")
	return nil
}
