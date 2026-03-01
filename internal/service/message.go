package service

import (
	"task-5/internal/model"
	"task-5/internal/repository"
)

type MessageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) GetAllMsgs() ([]model.Message, error) {
	return s.repo.FindAll()
}

func (s *MessageService) GetByID(id uint) (*model.Message, error) {
	return s.repo.FindByID(id)
}

func (s *MessageService) CreateMessage(msg *model.Message) error {
	return s.repo.Create(msg)
}

func (s *MessageService) DeleteMessage(id uint) error {
	return s.repo.Delete(id)
}
