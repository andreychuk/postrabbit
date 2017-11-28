package main

import (
	"github.com/caarlos0/env"
	"github.com/lib/pq"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Config struct {
	CHANNEL_LIST          string `env:"CHANNEL_LIST"`
	POSTGRES_URL          string `env:"POSTGRES_URL"`
	RABBITMQ_URL          string `env:"RABBITMQ_URL"`
	DEFAULT_EXCHANGE_NAME string `env:"DEFAULT_EXCHANGE_NAME" envDefault:""`
	DELAY_EXCHANGE_NAME   string `env:"DELAY_EXCHANGE_NAME" envDefault:"delay"`
}

var AmqpConnection *amqp.Connection
var Listener *pq.Listener
var Conf Config

func init() {
	env.Parse(&Conf)
	Listener = initNewListener(Conf)
	AmqpConnection = initRabbitConn(Conf)
}

func initRabbitConn(conf Config) *amqp.Connection {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		conn, err := amqp.Dial(conf.RABBITMQ_URL)
		if err != nil {
			log.Printf("[RABBIT CONNECTION] %s\n", err.Error())
			log.Println(
				"[RABBIT CONNECTION] node will only be able to respond to local connections")
			log.Println("[RABBIT CONNECTION] trying to reconnect in 5 seconds...")
			continue
		}
		return conn
	}
}

func errorReporter(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Printf("[ERROR REPORTER] %s\n", err.Error())
	}
}

func initNewListener(conf Config) *pq.Listener {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		listener := pq.NewListener(conf.POSTGRES_URL, 10*time.Second, time.Minute, errorReporter)

		channels := parseChannelList(conf.CHANNEL_LIST)

		for _, ch := range channels {
			err := listener.Listen(ch)
			if err != nil {
				log.Fatalf("[LISTENER ERROR] %s\n", err.Error())
			}
			log.Printf("[LISTENER] Start to listen channel %s\n", ch)
		}
		err := listener.Ping()
		if err != nil {
			log.Printf("[PQ CONNECTION] %s\n", err.Error())
			log.Println("[PQ CONNECTION] node will only be able to respond to local connections")
			log.Println("[RABBIT CONNECTION] trying to reconnect in 5 seconds...")
			continue
		}
		return listener
	}
}

func main() {
	run()
}
