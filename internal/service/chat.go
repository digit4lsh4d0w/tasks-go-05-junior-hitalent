package service

import (
	"errors"

	"task-5/internal/model"
)

type ChatRepository interface {
	FindAll() ([]model.Chat, error)
	FindByID(id uint) (*model.Chat, error)
	Create(chat *model.Chat) error
	Update(chat *model.Chat) error
	Delete(id uint) error
}

type chatService struct {
	repo ChatRepository
}

func NewChatService(repo ChatRepository) *chatService {
	return &chatService{repo: repo}
}

func (s *chatService) GetAllChats() ([]model.Chat, error) {
	return s.repo.FindAll()
}

func (s *chatService) GetChat(id uint) (*model.Chat, error) {
	return s.repo.FindByID(id)
}

func (s *chatService) CreateChat(chat *model.Chat) error {
	if chat.Title == "" {
		return errors.New("chat title is required")
	}
	return s.repo.Create(chat)
}

func (s *chatService) DeleteChat(id uint) error {
	return s.repo.Delete(id)
}
