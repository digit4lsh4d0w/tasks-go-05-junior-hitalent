package gorm

import (
	"task-5/internal/model"
	"task-5/internal/repository"

	"gorm.io/gorm"
)

type chatRepositoryGorm struct {
	db *gorm.DB
}

func NewChatRepositoryGorm(db *gorm.DB) repository.ChatRepository {
	return &chatRepositoryGorm{db: db}
}

func (r *chatRepositoryGorm) FindAll() ([]model.Chat, error) {
	var chats []model.Chat
	result := r.db.Find(&chats)
	return chats, result.Error
}

func (r *chatRepositoryGorm) FindByID(id uint) (*model.Chat, error) {
	var chat model.Chat
	result := r.db.Preload("Messages").First(&chat, id)
	return &chat, result.Error
}

func (r *chatRepositoryGorm) Create(chat *model.Chat) error {
	return r.db.Create(chat).Error
}

func (r *chatRepositoryGorm) Update(chat *model.Chat) error {
	return r.db.Save(chat).Error
}

func (r *chatRepositoryGorm) Delete(id uint) error {
	return r.db.Delete(&model.Chat{}, id).Error
}
