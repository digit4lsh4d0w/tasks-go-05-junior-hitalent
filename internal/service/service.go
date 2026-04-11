package service

import (
	"task-5/internal/model"
)

type ChatRepository interface {
	Create(chat *model.Chat) error
	FindByIDWithMessages(id uint, limit int) (*model.Chat, error)
	Delete(id uint) error
	CreateMessage(msg *model.Message) error
}

type chatService struct {
	chatRepo ChatRepository
}

func NewChatService(chatRepo ChatRepository) *chatService {
	return &chatService{chatRepo}
}

func (s *chatService) CreateChat(chat *model.Chat) error {
	return s.chatRepo.Create(chat)
}

func (s *chatService) GetChatWithMessages(id uint, limit int) (*model.Chat, error) {
	return s.chatRepo.FindByIDWithMessages(id, limit)
}

func (s *chatService) DeleteChat(id uint) error {
	return s.chatRepo.Delete(id)
}

func (s *chatService) CreateMessage(message *model.Message) error {
	return s.chatRepo.CreateMessage(message)
}
