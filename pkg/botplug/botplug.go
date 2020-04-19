package botplug

import (
	"time"
)

type BotPlugin interface {
	ReceiveMessage(input *MessageInput) (output *MessageOutput)
}

type MessageInput struct {
	Timestamp time.Time
	Source    *Source
	Messages  []string
}

type Source struct {
	Type    string
	UserID  string
	GroupID string
}

type MessageOutput struct {
	Queue []interface{}
}
