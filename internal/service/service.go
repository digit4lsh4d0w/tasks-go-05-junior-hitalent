package service

import (
	"errors"

	"task-5/internal/model"
)

type ChatRepository interface {
	FindAll() ([]model.Chat, error)
	FindByID(id uint) (*model.Chat, error)
	FindByIDWithMessages(id uint, limit int) (*model.Chat, error)
	Create(chat *model.Chat) error
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
	if chat.Title == "" {
		return errors.New("chat title is required")
	}
	return s.chatRepo.Create(chat)
}

func (s *chatService) GetAllChats() ([]model.Chat, error) {
	return s.chatRepo.FindAll()
}

func (s *chatService) GetChat(id uint) (*model.Chat, error) {
	return s.chatRepo.FindByID(id)
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
