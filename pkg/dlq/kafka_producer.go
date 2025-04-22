package dlq

import (
	"context"
	"encoding/json"
	"log"

	kafkaService "github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	Produce(ctx context.Context, topic, key string, value any) error
}

type SimpleKafkaProducer struct {
	service *kafkaService.KafkaService
}

func NewKafkaProducer(brokers []string, topic string) *SimpleKafkaProducer {
	return &SimpleKafkaProducer{
		service: kafkaService.NewProducer(brokers, topic),
	}
}

func (p *SimpleKafkaProducer) Produce(ctx context.Context, topic, key string, value any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		log.Printf("[DLQ] 🔴 Ошибка сериализации: %v", err)
		return err
	}

	return p.service.ProduceMessage(ctx, kafka.Message{
		Key:   []byte(key),
		Value: bytes,
		// Topic указывать нельзя, он уже зашит в Writer при создании
	})
}
