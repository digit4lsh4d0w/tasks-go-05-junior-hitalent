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
	if err == nil {
		chat.ID = dao.ID
		chat.CreatedAt = dao.CreatedAt
	}
	return err
}

func (r *chatRepository) FindAll() ([]model.Chat, error) {
	var daoChats []gormChat
	result := r.db.Find(&daoChats)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	}

	chats := make([]model.Chat, len(daoChats))
	for i, c := range daoChats {
		chats[i] = *toModelChat(&c)
	}
	return chats, result.Error
}

func (r *chatRepository) FindByID(id uint) (*model.Chat, error) {
	var c gormChat
	result := r.db.First(&c, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	}

	return toModelChat(&c), result.Error
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
	if err == nil {
		msg.ID = dao.ID
		msg.CreatedAt = dao.CreatedAt
	}

	return err
}
