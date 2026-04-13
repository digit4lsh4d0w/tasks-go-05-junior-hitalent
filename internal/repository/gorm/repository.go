package gorm

import (
	"context"
	"errors"

	"task-5/internal/model"

	"gorm.io/gorm"
)

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *chatRepository {
	return &chatRepository{db}
}

func (r *chatRepository) Create(ctx context.Context, chat *model.Chat) error {
	dao := toDAOChat(chat)
	if err := r.db.Create(dao).Error; err != nil {
		return err
	}

	chat.ID = dao.ID
	chat.CreatedAt = dao.CreatedAt

	return nil
}

func (r *chatRepository) FindByIDWithMessages(ctx context.Context, id uint, limit int) (*model.Chat, error) {
	var c gormChat
	err := r.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(limit)
	}).First(&c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound
		}

		return nil, err
	}

	return toModelChat(&c), nil
}

func (r *chatRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.Delete(&gormChat{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (r *chatRepository) CreateMessage(ctx context.Context, msg *model.Message) error {
	dao := toDAOMessage(msg)
	if err := r.db.Create(&dao).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return model.ErrNotFound
		}

		return err
	}

	msg.ID = dao.ID
	msg.CreatedAt = dao.CreatedAt

	return nil
}
