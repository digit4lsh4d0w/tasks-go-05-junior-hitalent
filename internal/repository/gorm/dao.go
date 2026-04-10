package gorm

import (
	"task-5/internal/model"

	"gorm.io/gorm"
)

type gormChat struct {
	gorm.Model
	Title    string        `gorm:"not null;unique"`
	Messages []gormMessage `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
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
	}

	if len(c.Messages) > 0 {
		chat.Messages = make([]model.Message, len(c.Messages))
		for i, m := range c.Messages {
			chat.Messages[i] = *toModelMessage(&m)
		}
	}

	return chat
}

func toDAOChat(m *model.Chat) *gormChat {
	if m == nil {
		return nil
	}
	c := &gormChat{}
	c.ID = m.ID
	c.Title = m.Title
	c.CreatedAt = m.CreatedAt
	return c
}

type gormMessage struct {
	gorm.Model
	ChatID uint   `gorm:"not null;index"`
	Text   string `gorm:"not null"`
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
	}
}

func toDAOMessage(m *model.Message) *gormMessage {
	if m == nil {
		return nil
	}
	msg := &gormMessage{}
	msg.ID = m.ID
	msg.ChatID = m.ChatID
	msg.Text = m.Text
	msg.CreatedAt = m.CreatedAt
	return msg
}
