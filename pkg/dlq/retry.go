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

func (r *KafkaRetrier) Retry(ctx context.Context, n Notification) error {
	return r.Producer.Produce(ctx, r.Topic, n.ID, n) // ✅ передаём весь объект
}
