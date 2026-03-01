package gorm

import (
	"task-5/internal/model"
	"task-5/internal/repository"

	"gorm.io/gorm"
)

type messageRepositoryGorm struct {
	db *gorm.DB
}

func NewMessageRepositoryGorm(db *gorm.DB) repository.MessageRepository {
	return &messageRepositoryGorm{db: db}
}

func (r *messageRepositoryGorm) FindAll() ([]model.Message, error) {
	var msgs []model.Message
	result := r.db.Find(&msgs)
	return msgs, result.Error
}

func (r *messageRepositoryGorm) FindByID(id uint) (*model.Message, error) {
	var msg model.Message
	result := r.db.Find(&msg, id)
	return &msg, result.Error
}

func (r *messageRepositoryGorm) Create(msg *model.Message) error {
	return r.db.Create(msg).Error
}

func (r *messageRepositoryGorm) Update(msg *model.Message) error {
	return r.db.Save(msg).Error
}

func (r *messageRepositoryGorm) Delete(id uint) error {
	return r.db.Delete(&model.Message{}, id).Error
}
