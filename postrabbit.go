package main

import (
	"github.com/caarlos0/env"
)

type Config struct {
	CHANNEL_LIST          string `env:"CHANNEL_LIST"`
	POSTGRES_URL          string `env:"POSTGRES_URL"`
	RABBITMQ_URL          string `env:"RABBITMQ_URL"`
	DEFAULT_EXCHANGE_NAME string `end:"DEFAULT_EXCHANGE_NAME" envDefault:"default"`
	DELAY_EXCHANGE_NAME   string `end:"DELAY_EXCHANGE_NAME" envDefault:"delay"`
}

func main() {
	conf := Config{}
	env.Parse(&conf)

	run(conf)
}
