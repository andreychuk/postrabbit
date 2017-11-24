package main

import (
	"github.com/caarlos0/env"
)

type Config struct {
	CHANNEL_LIST string `env:"CHANNEL_LIST"`
	POSTGRES_URL string `env:"POSTGRES_URL"`
	RABBITMQ_URL string `env:"RABBITMQ_URL"`
}

func main() {
	conf := Config{}
	env.Parse(&conf)

	run(conf)
}