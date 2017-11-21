package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/caarlos0/env"
	"postrabbit/config"
)

var (
	app          = kingpin.New("postrabbit", "A PostgreSQL/RabbitMQ Bridge")
	runcommand   = app.Command("run", "run the listener")
)

func main() {
	conf := config.Config{}
	env.Parse(&conf)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case runcommand.FullCommand():run(conf)
	}
}