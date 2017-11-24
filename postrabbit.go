package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/caarlos0/env"
)

type Config struct {
	CHANNEL_LIST string `env:"CHANNEL_LIST"`
	POSTGRES_URL string `env:"POSTGRES_URL"`
	RABBITMQ_URL string `env:"RABBITMQ_URL"`
}

var (
	app          = kingpin.New("postrabbit", "A PostgreSQL/RabbitMQ Bridge")
	runcommand   = app.Command("run", "run the listener")
)

func main() {
	conf := Config{}
	env.Parse(&conf)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case runcommand.FullCommand():run(conf)
	}
}