package main

import (
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/streadway/amqp"
	"strings"
)

func run() {
	rabbitChannel := make(chan pq.Notification, 100)

	go func() {
		defer AmqpConnection.Close()
		ch, err := AmqpConnection.Channel()

		if err != nil {
			log.Fatalf("[CHANNEL ERROR] %s\n", err.Error())
		}
		defer ch.Close()

		for {
			var msg Message
			notification := <-rabbitChannel
			err = ffjson.Unmarshal([]byte(notification.Extra), &msg)
			msg.Channel = notification.Channel
			msg.Data = notification.Extra
			msg.Exchange = Conf.DEFAULT_EXCHANGE_NAME
			if err != nil {
				log.Printf("[JOSN ERROR] %s\n", err.Error())
			} else {
				headers := make(map[string]interface{})
				if msg.isDelay() == true {
					// Delay messages
					if msg.getDelay() <= 0 {
						continue
					}
					headers["x-delay"] = msg.getDelay()
					msg.Exchange = Conf.DELAY_EXCHANGE_NAME
				}
				err := ch.Publish(msg.getExchange(), msg.getChannel(), false, false, amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(msg.getData()),
					Headers:     headers,
				})
				if err != nil {
					log.Fatalf("[PUBLISH ERROR] %s\n", err.Error())
				}
				log.Printf(msg.toString())
			}

		}

	}()

	for {
		select {
		case notification := <-Listener.Notify:
			rabbitChannel <- *notification
		case <-time.After(90 * time.Second):
			go func() {
				err := Listener.Ping()
				if err != nil {
					log.Fatalf("[LISTENER PING ERROR] %s\n", err.Error())
				}
			}()
		}
	}

}

func parseChannelList(list string) []string {
	channels := strings.Split(list, ",")
	return channels
}
