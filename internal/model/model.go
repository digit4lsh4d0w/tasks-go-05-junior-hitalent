package model

import (
	"strings"
	"time"
)

const (
	chatTitleMaxLen   = 200
	messageTextMaxLen = 5000
)

type Chat struct {
	ID        uint
	Title     string
	Messages  []Message
	CreatedAt time.Time
}

func NewChat(title string) (*Chat, error) {
	trimmedTitle := strings.TrimSpace(title)

	if trimmedTitle == "" {
		return nil, ErrChatTitleIsEmpty
	}
	if len([]rune(trimmedTitle)) > chatTitleMaxLen {
		return nil, ErrChatTitleTooLong
	}

	return &Chat{Title: trimmedTitle}, nil
}

type Message struct {
	ID        uint
	ChatID    uint
	Text      string
	CreatedAt time.Time
}

func NewMessage(text string) (*Message, error) {
	if text == "" {
		return nil, ErrMessageTextIsEmpty
	}
	if len([]rune(text)) > messageTextMaxLen {
		return nil, ErrMessageTextTooLong
	}

	return &Message{Text: text}, nil
}
