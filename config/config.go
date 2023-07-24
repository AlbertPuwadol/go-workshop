package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBURI    string `env:"MONGODB_URI"`
	RabbitMQURI   string `env:"RABBITMQ_URI"`
	IntervalQueue string `env:"INTERVAL_QUEUE"`
	ResultQueue   string `env:"RESULT_QUEUE"`
	ApiURL        string `env:"API_URL"`
}

func NewConfig() Config {
	godotenv.Load()
	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}
	return config
}
