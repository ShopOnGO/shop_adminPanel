package dlq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	Produce(ctx context.Context, topic, key string, value any) error
}

type SimpleKafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) *SimpleKafkaProducer {
	return &SimpleKafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *SimpleKafkaProducer) Produce(ctx context.Context, topic, key string, value any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		log.Printf("[DLQ] 🔴 Ошибка сериализации: %v", err)
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: bytes,
	})
}
