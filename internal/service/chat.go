package service

import (
	"errors"

	"task-5/internal/model"
	"task-5/internal/repository"
)

type ChatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) GetAllChats() ([]model.Chat, error) {
	return s.repo.FindAll()
}

func (s *ChatService) GetChat(id uint) (*model.Chat, error) {
	return s.repo.FindByID(id)
}

func (s *ChatService) CreateChat(chat *model.Chat) error {
	if chat.Title == "" {
		return errors.New("chat title is required")
	}
	return s.repo.Create(chat)
}

func (s *ChatService) DeleteChat(id uint) error {
	return s.repo.Delete(id)
}
