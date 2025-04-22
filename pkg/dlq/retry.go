package dlq

import (
	"context"
)

type KafkaRetrier struct {
	Producer KafkaProducer // Интерфейс обёртки над Kafka producer
	Topic    string
}

func NewKafkaRetrier(producer KafkaProducer, topic string) *KafkaRetrier {
	return &KafkaRetrier{
		Producer: producer,
		Topic:    topic,
	}
}

func (r *KafkaRetrier) RetryWithKey(ctx context.Context, key string, n Notification) error {
	return r.Producer.Produce(ctx, r.Topic, key, n) // ✅ передаём весь объект
}
