package db

import (
	"time"
)

type Message struct {
	Author string
	Data   string
	Time   time.Time
}

type FrontEndMessage struct {
	Author string
	Data   string
	Time   string
}

func GetMessages() []FrontEndMessage {
	message := []FrontEndMessage{{
		Author: "Test",
		Data:   "Lorem Ipsum",
		Time:   time.Now().Format(time.DateTime),
	}}
	return message
}
