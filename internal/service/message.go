package service

import (
	"task-5/internal/model"
)

type MessageRepository interface {
	FindAll() ([]model.Message, error)
	FindByID(id uint) (*model.Message, error)
	Create(msg *model.Message) error
	Update(msg *model.Message) error
	Delete(id uint) error
}

type messageService struct {
	repo MessageRepository
}

func NewMessageService(repo MessageRepository) *messageService {
	return &messageService{repo: repo}
}

func (s *messageService) GetAllMsgs() ([]model.Message, error) {
	return s.repo.FindAll()
}

func (s *messageService) GetByID(id uint) (*model.Message, error) {
	return s.repo.FindByID(id)
}

func (s *messageService) CreateMessage(msg *model.Message) error {
	return s.repo.Create(msg)
}

func (s *messageService) DeleteMessage(id uint) error {
	return s.repo.Delete(id)
}
