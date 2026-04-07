package service

import (
	"errors"

	"task-5/internal/model"
)

type ChatRepository interface {
	FindAll() ([]model.Chat, error)
	FindByID(id uint) (*model.Chat, error)
	Create(chat *model.Chat) error
	Delete(id uint) error
}

type MessageRepository interface {
	FindAll() ([]model.Message, error)
	FindByID(id uint) (*model.Message, error)
	Create(msg *model.Message) error
	Delete(id uint) error
}

type chatService struct {
	chatRepo ChatRepository
	msgRepo  MessageRepository
}

func NewChatService(cr ChatRepository, mr MessageRepository) *chatService {
	return &chatService{
		chatRepo: cr,
		msgRepo:  mr,
	}
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
	return s.chatRepo.FindByID(id)
}

func (s *chatService) DeleteChat(id uint) error {
	return s.chatRepo.Delete(id)
}

func (s *chatService) CreateMessage(message *model.Message) error {
	return s.msgRepo.Create(message)
}
