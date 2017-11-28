package main

import "fmt"

type (
	Message struct {
		Channel  string
		Delay    int64
		Data     string
		Exchange string
	}
)

func (msg *Message) getChannel() string {
	return msg.Channel
}

func (msg *Message) isDelay() bool {
	if msg.Delay > 0 {
		return true
	} else {
		return false
	}
}

func (msg *Message) getDelay() int64 {
	if msg.Delay > 0 {
		return msg.Delay
	} else {
		return 0
	}
}

func (msg *Message) getData() string {
	return msg.Data
}

func (msg *Message) toString() string {
	if msg.isDelay() == true {
		return fmt.Sprintf(
			"[DELAY MESSAGE] ch: %s, scheduled to: %d ms, data: %s\n",
			msg.getChannel(),
			msg.getDelay(),
			msg.getData(),
		)
	} else {
		return fmt.Sprintf(
			"[SIMPLE MESSAGE] ch: %s, data: %s\n",
			msg.getChannel(),
			msg.getData(),
		)
	}

}
func (msg *Message) getExchange() string {
	return msg.Exchange
}
