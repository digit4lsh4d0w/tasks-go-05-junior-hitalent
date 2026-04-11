package gorm

import (
	"time"

	"task-5/internal/model"
)

type gormChat struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []gormMessage `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
}

func (gormChat) TableName() string {
	return "chats"
}

func toModelChat(c *gormChat) *model.Chat {
	if c == nil {
		return nil
	}

	chat := &model.Chat{
		ID:        c.ID,
		Title:     c.Title,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	chat.Messages = make([]model.Message, len(c.Messages))
	for i, m := range c.Messages {
		chat.Messages[i] = *toModelMessage(&m)
	}

	return chat
}

func toDAOChat(m *model.Chat) *gormChat {
	if m == nil {
		return nil
	}
	return &gormChat{
		ID:        m.ID,
		Title:     m.Title,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type gormMessage struct {
	ID        uint   `gorm:"primaryKey"`
	ChatID    uint   `gorm:"not null;index"`
	Text      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (gormMessage) TableName() string {
	return "messages"
}

func toModelMessage(m *gormMessage) *model.Message {
	if m == nil {
		return nil
	}
	return &model.Message{
		ID:        m.ID,
		ChatID:    m.ChatID,
		Text:      m.Text,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func toDAOMessage(m *model.Message) *gormMessage {
	if m == nil {
		return nil
	}
	return &gormMessage{
		ID:        m.ID,
		ChatID:    m.ChatID,
		Text:      m.Text,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
