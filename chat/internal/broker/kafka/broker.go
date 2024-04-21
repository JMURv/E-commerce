package kafka

import (
	"github.com/IBM/sarama"
	conf "github.com/JMURv/e-commerce/chat/pkg/config"
	"log"
	"strconv"
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

func (b *Broker) NewMessageNotification(msgID uint64, msg []byte) error {
	b.producer.Input() <- &sarama.ProducerMessage{
		Topic: b.topic,
		Key:   sarama.StringEncoder(strconv.FormatUint(msgID, 10)),
		Value: sarama.ByteEncoder(msg),
	}
	log.Println("Message sent successfully!")
	return nil
}
