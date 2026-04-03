package service

import (
	"errors"

	"task-5/internal/model"
	"task-5/internal/repository"
)

type chatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) ChatService {
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
