package kafka

import (
	"github.com/IBM/sarama"
	conf "github.com/JMURv/e-commerce/chat/pkg/config"
	"log"
	"strconv"
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

func (b *Broker) NewMessageNotification(msgID uint64, msg []byte) error {
	b.producer.Input() <- &sarama.ProducerMessage{
		Topic: b.cfg.Kafka.NotificationTopic,
		Key:   sarama.StringEncoder(strconv.FormatUint(msgID, 10)),
		Value: sarama.ByteEncoder(msg),
	}
	log.Println("Message sent successfully!")
	return nil
}
