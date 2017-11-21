package config

type (
	Config struct {
		CHANNEL_LIST string `env:"CHANNEL_LIST"`
		POSTGRES_URL string `env:"POSTGRES_URL"`
		RABBITMQ_URL string `env:"RABBITMQ_URL"`
	}
)
