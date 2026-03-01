package repository

import (
	"task-5/internal/model"
)

type ChatRepository interface {
	FindAll() ([]model.Chat, error)
	FindByID(id uint) (*model.Chat, error)
	Create(chat *model.Chat) error
	Update(chat *model.Chat) error
	Delete(id uint) error
}

type MessageRepository interface {
	FindAll() ([]model.Message, error)
	FindByID(id uint) (*model.Message, error)
	Create(msg *model.Message) error
	Update(msg *model.Message) error
	Delete(id uint) error
}
