package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"task-5/internal/log"
	"task-5/internal/model"

	"github.com/go-playground/validator"
)

const (
	limitDefault = 5
	limitMax     = 20
)

type ChatService interface {
	CreateChat(chat *model.Chat) error
	GetChat(id uint) (*model.Chat, error)
	GetChatWithMessages(id uint, limit int) (*model.Chat, error)
	DeleteChat(id uint) error
	CreateMessage(msg *model.Message) error
}

type chatHandler struct {
	baseHandler
	service   ChatService
	validator *validator.Validate
}

type CreateChatRequest struct {
	Title string `json:"title" validate:"required,min=1,max=200"`
}

type SendMessageRequest struct {
	Text string `json:"text" validate:"required,min=1,max=5000"`
}

func NewChatHandler(s ChatService, v *validator.Validate, l log.Logger) *chatHandler {
	return &chatHandler{
		baseHandler: NewBaseHandler(l),
		service:     s,
		validator:   v,
	}
}

func parseChatID(r *http.Request) (uint, error) {
	chatIDStr := r.PathValue("chat_id")
	chatID, err := strconv.ParseUint(chatIDStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(chatID), nil
}

func parseLimit(r *http.Request) int {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return limitDefault
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return limitDefault
	}

	if limit <= 0 {
		return limitDefault
	}

	if limit > limitMax {
		return limitMax
	}

	return limit
}

func (h *chatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	// Ограничение тела запроса до 1 МебиБайта
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("failed to decode json body", "error", err)
		h.respondError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.log.Warn("failed to validate request", "error", err)
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	chat, err := model.NewChat(req.Title)
	if err != nil {
		h.log.Error("failed to convert chat DTO to model", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to create chat")
		return
	}

	if err := h.service.CreateChat(chat); err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			h.log.Error("chat already exists", "error", err)
			h.respondError(w, http.StatusConflict, "chat already exists")
			return
		}

		h.log.Error("failed to create chat", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to create chat")
		return
	}

	h.respondJSON(w, http.StatusCreated, chat)
}

func (h *chatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	chatID, err := parseChatID(r)
	if err != nil {
		h.log.Warn("failed to parse chat id", "error", err)
		h.respondError(w, http.StatusBadRequest, "invalid chat id")
		return
	}

	if err = h.service.DeleteChat(chatID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.log.Warn("chat not found", "error", err, "chat_id", chatID)
			h.respondError(w, http.StatusNotFound, "chat not found")
			return
		}

		h.log.Error("failed to delete chat", "error", err, "chat_id", chatID)
		h.respondError(w, http.StatusInternalServerError, "failed to delete chat")
		return
	}

	h.respondSuccess(w, http.StatusOK, "chat deleted successfully")
}

func (h *chatHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	chatID, err := parseChatID(r)
	if err != nil {
		h.log.Warn("failed to parse chat id", "error", err)
		h.respondError(w, http.StatusBadRequest, "invalid chat id")
		return
	}

	// Ограничение тела запроса до 2 МебиБайт
	r.Body = http.MaxBytesReader(w, r.Body, 2*1<<20)
	defer r.Body.Close()

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("failed to decode json body", "error", err)
		h.respondError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.log.Warn("invalid request body", "error", err)
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	message := &model.Message{
		ChatID: chatID,
		Text:   req.Text,
	}

	if err := h.service.CreateMessage(message); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.log.Warn("chat not found", "error", err)
			h.respondError(w, http.StatusNotFound, "chat not found")
			return
		}

		h.log.Error("failed to create message", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to create message")
		return
	}

	h.respondJSON(w, http.StatusCreated, message)
}

func (h *chatHandler) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	chatID, err := parseChatID(r)
	if err != nil {
		h.log.Warn("failed to parse chat id", "error", err)
		h.respondError(w, http.StatusBadRequest, "invalid chat id")
		return
	}

	limit := parseLimit(r)

	chat, err := h.service.GetChatWithMessages(chatID, limit)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.log.Warn("chat not found", "error", err, "chat_id", chatID)
			h.respondError(w, http.StatusNotFound, "chat not found")
			return
		}

		h.log.Error("failed to get chat with messages", "error", err, "chat_id", chatID)
		h.respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.respondJSON(w, http.StatusOK, chat)
}
