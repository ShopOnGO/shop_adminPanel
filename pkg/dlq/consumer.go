package dlq

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader   *kafka.Reader
	dispatch func([]byte)
}

func NewConsumer(brokers []string, topic, groupID string, dispatch func([]byte)) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 1,
			MaxBytes: 10e6,
		}),
		dispatch: dispatch,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	defer c.reader.Close()

	log.Println("[DLQ] ⏳ Стартуем потребление из Kafka...")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			// Завершение по отмене контекста
			if ctx.Err() != nil {
				log.Println("[DLQ] 🛑 Консьюмер остановлен")
				return
			}
			log.Printf("[DLQ] ❌ Ошибка чтения сообщения: %v", err)
			continue
		}

		log.Printf("[DLQ] 📩 Сообщение: key=%s, value=%s", string(m.Key), string(m.Value))
		c.dispatch(m.Value)
	}
}

// Вариант запуска с graceful shutdown
func RunWithGracefulShutdown(consumer *Consumer) {
	ctx, cancel := context.WithCancel(context.Background())

	go consumer.Start(ctx)

	// Ожидаем сигнал завершения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("[DLQ] 🔄 Завершаем...")
	cancel()
	time.Sleep(2 * time.Second) // Дать время на завершение
}
