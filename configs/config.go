package configs

import (
	"os"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"

	"github.com/ShopOnGO/ShopOnGO/configs"
	"github.com/joho/godotenv"
)

type Config struct {
	Db           DbConfig
	Dlq          DlqConfig
	LogLevel     logger.LogLevel
	FileLogLevel logger.LogLevel
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
	// logger
	logLevelStr := os.Getenv("ADMIN_SERVICE_LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = "INFO"
	}
	LogLevel := configs.ParseLogLevel(logLevelStr)
	fileLogLevelStr := os.Getenv("ADMIN_SERVICE_FILE_LOG_LEVEL")
	if fileLogLevelStr == "" {
		fileLogLevelStr = "INFO"
	}
	FileLogLevel := configs.ParseLogLevel(fileLogLevelStr)

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Dlq: DlqConfig{
			Broker:        os.Getenv("KAFKA_BROKER"),
			ConsumerTopic: os.Getenv("KAFKA_CONSUMER"),
			ProducerTopic: os.Getenv("KAFKA_PRODUCER"),
		},
		LogLevel:     LogLevel,
		FileLogLevel: FileLogLevel,
	}
}
