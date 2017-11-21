package main

import (
	"log"
	"time"

	pq "github.com/lib/pq"
	"github.com/streadway/amqp"
	"postrabbit/config"
	"strings"
	"github.com/pquerna/ffjson/ffjson"
)

func errorReporter(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Print(err)
	}
}

func run(conf config.Config) {
	listener := pq.NewListener(conf.POSTGRES_URL, 10*time.Second, time.Minute, errorReporter)
	channels := parseChannelList(conf.CHANNEL_LIST)

	for _, ch := range channels {
		err := listener.Listen(ch)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[LISTENER] Start to listen channel %s\n", ch)
	}

	rabbitChannel := make(chan string, 100)

	go func() {
		conn, err := amqp.Dial(conf.RABBITMQ_URL)

		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		ch, err := conn.Channel()

		if err != nil {
			log.Fatal(err)
		}
		defer ch.Close()

		for {
			var msg Message
			payload := <-rabbitChannel
			err := ffjson.Unmarshal([]byte(payload), &msg)
			if err != nil {
				log.Println(err)
			} else {
				headers := make(map[string]interface{})
				if msg.isDelay() == true {
					// Delay messages
					headers["x-delay"] = msg.getDelay()
				}
				err := ch.Publish(msg.getChannel(), "", false, false, amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(msg.getData()),
					Headers:	 headers,
				})
				if err != nil {
					log.Fatal(err)
				}
				log.Printf(msg.toString())
			}

		}

	}()

	for {
		select {
		case notification := <-listener.Notify:
			rabbitChannel <- notification.Extra
		case <-time.After(90 * time.Second):
			go func() {
				err := listener.Ping()
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	}

}

func parseChannelList(list string) ([]string) {
	channels := strings.Split(list, ",")
	return channels
}