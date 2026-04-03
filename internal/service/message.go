package service

import (
	"task-5/internal/model"
	"task-5/internal/repository"
)

type messageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) MessageService {
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
