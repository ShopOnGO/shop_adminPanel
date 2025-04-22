package dlq

import (
	"admin/configs"
	"context"
	"encoding/json"
	"log"
)

func StartDLQProcessor(conf *configs.Config) {
	retrier := NewKafkaRetrier(
		NewKafkaProducer([]string{conf.Dlq.Broker}, conf.Dlq.ProducerTopic),
		conf.Dlq.ProducerTopic,
	)

	log.Printf("[DLQ] ⚙️ Broker: %s, Topic: %s", conf.Dlq.Broker, conf.Dlq.ConsumerTopic)

	consumer := NewConsumer(
		[]string{conf.Dlq.Broker},
		conf.Dlq.ConsumerTopic,
		"dlq-processor-group",
		"dlq-processor-client", // clientID
	)

	handler := func(key, value []byte) {
		var n Notification
		if err := json.Unmarshal(value, &n); err != nil {
			log.Printf("[DLQ] 🚫 Не удалось распарсить сообщение: %v", err)
			return
		}

		if ShouldRetry(n) {
			n.WasInDLQ = true
			log.Printf("[DLQ] 🔁 Повторная отправка: %s", n.ID)
			if err := retrier.RetryWithKey(context.Background(), string(key), n); err != nil {
				log.Printf("[DLQ] ❌ Ошибка ретрая: %v", err)
			}
		} else {
			log.Printf("[DLQ] ❎ Пропускаем: %s", n.ID)
		}
	}

	RunWithGracefulShutdown(consumer, handler)
}
