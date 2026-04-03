package handler

import (
	"net/http"

	"task-5/internal/service"
)

type ChatHandler struct {
	service *service.ChatService
}

func NewChatHandler(s *service.ChatService) *ChatHandler {
	return &ChatHandler{service: s}
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {}

func (h *ChatHandler) SendMsg(w http.ResponseWriter, r *http.Request) {}

func (h *ChatHandler) GetMsgs(w http.ResponseWriter, r *http.Request) {}

func (h *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {}
