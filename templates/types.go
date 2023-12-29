package templates

type Message struct {
	Author string
	Data   string
	Time   string
	Photo  string
}

type MessageList []Message
