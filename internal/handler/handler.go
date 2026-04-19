package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"task-5/internal/model"

	"github.com/go-playground/validator"
)

const (
	limitDefault = 5
	limitMax     = 20
)

type ChatService interface {
	CreateChat(ctx context.Context, chat *model.Chat) error
	GetChatWithMessages(ctx context.Context, id uint, limit int) (*model.Chat, error)
	DeleteChat(ctx context.Context, id uint) error
	CreateMessage(ctx context.Context, msg *model.Message) error
}

type chatHandler struct {
	baseHandler
	service   ChatService
	validator *validator.Validate
}

func NewChatHandler(service ChatService, validator *validator.Validate, logger *slog.Logger) *chatHandler {
	return &chatHandler{
		baseHandler: NewBaseHandler(logger),
		service:     service,
		validator:   validator,
	}
}

func (h *chatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	// Ограничение тела запроса до 1 МебиБайта
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WarnContext(r.Context(), "failed to decode json body",
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.WarnContext(r.Context(), "validation error",
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusBadRequest, "validation error")
		return
	}

	chat, err := model.NewChat(req.Title)
	if err != nil {
		h.logger.WarnContext(r.Context(), "invalid chat data",
			slog.String("chat_title", req.Title),
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusUnprocessableEntity, "invalid chat data")
		return
	}

	if err := h.service.CreateChat(r.Context(), chat); err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			h.logger.WarnContext(r.Context(), "chat already exists",
				slog.String("error", err.Error()),
			)
			h.respondError(r.Context(), w, http.StatusConflict, "chat already exists")
			return
		}

		h.logger.ErrorContext(r.Context(), "failed to create chat",
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusInternalServerError, "failed to create chat")
		return
	}

	h.respondJSON(r.Context(), w, http.StatusCreated, chat)
}

func (h *chatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.PathValue("chat_id")
	chatID, err := parseChatID(chatIDStr)
	if err != nil {
		h.logger.WarnContext(r.Context(), "failed to parse chat id",
			slog.String("chat_id", chatIDStr),
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusBadRequest, "invalid chat id")
		return
	}

	if err = h.service.DeleteChat(r.Context(), chatID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.logger.WarnContext(r.Context(), "chat not found",
				slog.Uint64("chat_id", uint64(chatID)),
				slog.String("error", err.Error()),
			)
			h.respondError(r.Context(), w, http.StatusNotFound, "chat not found")
			return
		}

		h.logger.ErrorContext(r.Context(), "failed to delete chat",
			slog.Uint64("chat_id", uint64(chatID)),
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusInternalServerError, "failed to delete chat")
		return
	}

	h.respondSuccess(r.Context(), w, http.StatusOK, "chat deleted successfully")
}

func (h *chatHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.PathValue("chat_id")
	chatID, err := parseChatID(chatIDStr)
	if err != nil {
		h.logger.WarnContext(r.Context(), "failed to parse chat id",
			slog.String("chat_id", chatIDStr),
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusBadRequest, "invalid chat id")
		return
	}

	// Ограничение тела запроса до 2 МебиБайт
	r.Body = http.MaxBytesReader(w, r.Body, 2*1<<20)
	defer r.Body.Close()

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WarnContext(r.Context(), "failed to decode json body",
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.WarnContext(r.Context(), "validation error",
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	message, err := model.NewMessage(chatID, req.Text)
	if err != nil {
		h.logger.WarnContext(r.Context(), "invalid message data",
			slog.String("chat_id", chatIDStr),
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusUnprocessableEntity, "invalid message data")
		return
	}

	if err := h.service.CreateMessage(r.Context(), message); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.logger.WarnContext(r.Context(), "chat not found",
				slog.String("error", err.Error()),
			)
			h.respondError(r.Context(), w, http.StatusNotFound, "chat not found")
			return
		}

		h.logger.ErrorContext(r.Context(), "failed to create message",
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusInternalServerError, "failed to create message")
		return
	}

	h.respondJSON(r.Context(), w, http.StatusCreated, message)
}

func (h *chatHandler) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.PathValue("chat_id")
	chatID, err := parseChatID(chatIDStr)
	if err != nil {
		h.logger.WarnContext(r.Context(), "failed to parse chat id",
			slog.String("chat_id", chatIDStr),
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusBadRequest, "invalid chat id")
		return
	}

	limit := parseLimit(r.URL.Query().Get("limit"))

	chat, err := h.service.GetChatWithMessages(r.Context(), chatID, limit)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.logger.WarnContext(r.Context(), "chat not found",
				slog.Uint64("chat_id", uint64(chatID)),
				slog.String("error", err.Error()),
			)
			h.respondError(r.Context(), w, http.StatusNotFound, "chat not found")
			return
		}

		h.logger.ErrorContext(r.Context(), "failed to get chat with messages",
			slog.Uint64("chat_id", uint64(chatID)),
			slog.String("error", err.Error()),
		)
		h.respondError(r.Context(), w, http.StatusInternalServerError, "internal error")
		return
	}

	h.respondJSON(r.Context(), w, http.StatusOK, chat)
}
