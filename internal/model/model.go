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

func NewMessage(chatID uint, text string) (*Message, error) {
	trimmedText := strings.TrimSpace(text)

	if trimmedText == "" {
		return nil, ErrMessageTextIsEmpty
	}
	if len([]rune(trimmedText)) > messageTextMaxLen {
		return nil, ErrMessageTextTooLong
	}

	return &Message{ChatID: chatID, Text: trimmedText}, nil
}
