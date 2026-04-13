package service

import (
	"context"
	"fmt"
	"log/slog"

	"task-5/internal/model"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) error
	FindByIDWithMessages(ctx context.Context, id uint, limit int) (*model.Chat, error)
	Delete(ctx context.Context, id uint) error
	CreateMessage(ctx context.Context, msg *model.Message) error
}

type chatService struct {
	chatRepo ChatRepository
	logger   *slog.Logger
}

func NewChatService(chatRepo ChatRepository, logger *slog.Logger) *chatService {
	return &chatService{chatRepo, logger}
}

func (s *chatService) CreateChat(ctx context.Context, chat *model.Chat) error {
	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return fmt.Errorf("repository.Create: %w", err)
	}

	s.logger.DebugContext(ctx, "chat created", slog.Uint64("chat_id", uint64(chat.ID)))
	return nil
}

func (s *chatService) GetChatWithMessages(ctx context.Context, id uint, limit int) (*model.Chat, error) {
	chat, err := s.chatRepo.FindByIDWithMessages(ctx, id, limit)
	if err != nil {
		return nil, fmt.Errorf("repository.FindByIDWithMessages: %w", err)
	}

	return chat, nil
}

func (s *chatService) DeleteChat(ctx context.Context, id uint) error {
	if err := s.chatRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}

	s.logger.DebugContext(ctx, "chat deleted", slog.Uint64("chat_id", uint64(id)))
	return nil
}

func (s *chatService) CreateMessage(ctx context.Context, message *model.Message) error {
	if err := s.chatRepo.CreateMessage(ctx, message); err != nil {
		return fmt.Errorf("repository.CreateMessage: %w", err)
	}

	s.logger.DebugContext(ctx, "message created", slog.Uint64("message_id", uint64(message.ID)))
	return nil
}
