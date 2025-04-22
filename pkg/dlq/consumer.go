package dlq

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	kafkaService "github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	kafka *kafkaService.KafkaService
}

func NewConsumer(brokers []string, topic, groupID, clientID string) *Consumer {
	return &Consumer{
		kafka: kafkaService.NewConsumer(brokers, topic, groupID, clientID),
	}
}

// Start запускает чтение сообщений из Kafka с обработчиком
func (c *Consumer) Start(ctx context.Context, handler func(key, value []byte)) {
	c.kafka.Consume(ctx, func(m kafka.Message) error {
		log.Printf("[DLQ] 📩 Сообщение: key=%s, value=%s", string(m.Key), string(m.Value))
		handler(m.Key, m.Value)
		return nil
	})
}

// Запуск с graceful shutdown
func RunWithGracefulShutdown(consumer *Consumer, handler func(key, value []byte)) {
	ctx, cancel := context.WithCancel(context.Background())

	go consumer.Start(ctx, handler)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("[DLQ] 🔄 Завершаем...")
	cancel()
	time.Sleep(2 * time.Second)
}
