package gorm

import (
	"errors"

	"task-5/internal/log"
	"task-5/internal/model"

	"gorm.io/gorm"
)

type chatRepository struct {
	db  *gorm.DB
	log log.Logger
}

func NewChatRepository(db *gorm.DB, log log.Logger) *chatRepository {
	log.Debug("Creating new chat repository", "type", "gorm")
	return &chatRepository{db, log}
}

func (r *chatRepository) Create(chat *model.Chat) error {
	dao := toDAOChat(chat)
	err := r.db.Create(dao).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return model.ErrAlreadyExists
		}
		return err
	}

	chat.ID = dao.ID
	chat.CreatedAt = dao.CreatedAt
	chat.UpdatedAt = dao.UpdatedAt

	return err
}

func (r *chatRepository) FindByIDWithMessages(id uint, limit int) (*model.Chat, error) {
	var c gormChat
	result := r.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(limit)
	}).First(&c, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound
		}
		return nil, result.Error
	}

	return toModelChat(&c), nil
}

func (r *chatRepository) Delete(id uint) error {
	result := r.db.Delete(&gormChat{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *chatRepository) CreateMessage(msg *model.Message) error {
	dao := toDAOMessage(msg)

	err := r.db.Create(&dao).Error
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return model.ErrNotFound
		}
		return err
	}

	msg.ID = dao.ID
	msg.CreatedAt = dao.CreatedAt
	msg.UpdatedAt = dao.UpdatedAt

	return nil
}
