package dlq

import (
	"admin/configs"
	"context"
	"encoding/json"
	"log"
)

func StartDLQProcessor(conf *configs.Config) {
	retrier := NewKafkaRetrier(
		NewKafkaProducer([]string{conf.Dlq.Broker}),
		conf.Dlq.ProducerTopic,
	)
	log.Printf("[DLQ] ⚙️ Broker: %s, Topic: %s", conf.Dlq.Broker, conf.Dlq.ConsumerTopic)

	consumer := NewConsumer(
		[]string{conf.Dlq.Broker},
		conf.Dlq.ConsumerTopic,
		"dlq-processor-group",
		func(msg []byte) {
			var n Notification
			if err := json.Unmarshal(msg, &n); err != nil {
				log.Printf("[DLQ] 🚫 Не удалось распарсить сообщение: %v", err)
				return
			}

			if ShouldRetry(n) {
				n.WasInDLQ = true
				log.Printf("[DLQ] 🔁 Повторная отправка: %s", n.ID)
				if err := retrier.Retry(context.Background(), n); err != nil {
					log.Printf("[DLQ] ❌ Ошибка ретрая: %v", err)
				}
			} else {
				log.Printf("[DLQ] ❎ Пропускаем: %s", n.ID)
			}
		},
	)

	RunWithGracefulShutdown(consumer)
}
