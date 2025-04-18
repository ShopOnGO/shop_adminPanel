package configs

import (
	"admin/pkg/logger"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db  DbConfig
	Dlq DlqConfig
}
type DlqConfig struct {
	Broker        string
	ConsumerTopic string
	ProducerTopic string
}
type DbConfig struct {
	Dsn string
}

func LoadConfig() *Config {
	err := godotenv.Load() //loading from .env
	if err != nil {
		logger.Error("Error loading .env file, using default config", err.Error())
	}

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Dlq: DlqConfig{
			Broker:        os.Getenv("KAFKA_BROKER"),
			ConsumerTopic: os.Getenv("KAFKA_CONSUMER"),
			ProducerTopic: os.Getenv("KAFKA_PRODUCER"),
		},
	}
}
