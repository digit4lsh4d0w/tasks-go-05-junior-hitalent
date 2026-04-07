package model

import "time"

type Chat struct {
	ID        uint
	Title     string
	Messages  []Message
	CreatedAt time.Time
}

type Message struct {
	ID        uint
	ChatID    uint
	Text      string
	CreatedAt time.Time
}
