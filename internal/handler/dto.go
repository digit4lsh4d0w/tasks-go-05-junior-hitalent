package handler

import (
	"time"

	"task-5/internal/model"
)

type CreateChatRequest struct {
	Title string `json:"title" validate:"required,min=1,max=200"`
}

type SendMessageRequest struct {
	Text string `json:"text" validate:"required,min=1,max=5000"`
}

type ChatResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func toChatResponse(chat model.Chat) ChatResponse {
	return ChatResponse{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
	}
}

type ChatDetailResponse struct {
	ChatResponse
	Messages []MessageResponse `json:"messages"`
}

func toChatDetailResponse(chat model.Chat) ChatDetailResponse {
	messages := make([]MessageResponse, 0, len(chat.Messages))
	for _, message := range chat.Messages {
		messages = append(messages, toMessageResponse(message))
	}

	return ChatDetailResponse{
		ChatResponse: toChatResponse(chat),
		Messages:     messages,
	}
}

type MessageResponse struct {
	ID        uint      `json:"id"`
	ChatID    uint      `json:"chat_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func toMessageResponse(message model.Message) MessageResponse {
	return MessageResponse{
		ID:        message.ID,
		ChatID:    message.ChatID,
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
	}
}
