package gorm

import (
	"task-5/internal/log"
	"task-5/internal/model"

	"gorm.io/gorm"
)

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

type messageRepository struct {
	db  *gorm.DB
	log log.Logger
}

func NewMessageRepository(db *gorm.DB, log log.Logger) *messageRepository {
	log.Debug("Creating new message repository", "type", "gorm")
	return &messageRepository{db, log}
}

func (r *messageRepository) Create(msg *model.Message) error {
	dao := toDAOMessage(msg)
	err := r.db.Create(&dao).Error
	if err == nil {
		msg.ID = dao.ID
		msg.CreatedAt = dao.CreatedAt
	}
	return err
}

func (r *messageRepository) FindAll() ([]model.Message, error) {
	var daoMsgs []gormMessage
	result := r.db.Find(&daoMsgs)
	if result.Error != nil {
		return nil, result.Error
	}
	msgs := make([]model.Message, len(daoMsgs))
	for i, m := range daoMsgs {
		msgs[i] = *toModelMessage(&m)
	}
	return msgs, nil
}

func (r *messageRepository) FindByID(id uint) (*model.Message, error) {
	var m gormMessage
	result := r.db.First(&m, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toModelMessage(&m), nil
}

func (r *messageRepository) Delete(id uint) error {
	return r.db.Delete(&gormMessage{}, id).Error
}
