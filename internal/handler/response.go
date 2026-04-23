package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"task-5/internal/model"
)

type baseHandler struct {
	logger *slog.Logger
}

func NewBaseHandler(l *slog.Logger) baseHandler {
	return baseHandler{logger: l}
}

func (h *baseHandler) respondJSON(ctx context.Context, w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if code != http.StatusNoContent && payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			h.logger.ErrorContext(ctx, "failed to encode response", "error", err.Error())
		}
	}
}

type successRespones struct {
	Success string `json:"success"`
}

func (h *baseHandler) respondSuccess(ctx context.Context, w http.ResponseWriter, code int, message string) {
	h.respondJSON(ctx, w, code, successRespones{message})
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *baseHandler) respondError(ctx context.Context, w http.ResponseWriter, code int, message string) {
	h.respondJSON(ctx, w, code, errorResponse{message})
}

// Обработка доменных ошибок
func (h *baseHandler) handleError(ctx context.Context, w http.ResponseWriter, err error, msg string, args ...any) {
	args = append(args, slog.String("error", err.Error()))

	switch {
	case errors.Is(err, model.ErrNotFound):
		h.logger.WarnContext(ctx, msg, args...)
		h.respondError(ctx, w, http.StatusNotFound, "not found")
	case errors.Is(err, model.ErrAlreadyExists):
		h.logger.WarnContext(ctx, msg, args...)
		h.respondError(ctx, w, http.StatusConflict, "already exists")
	default:
		h.logger.ErrorContext(ctx, msg, args...)
		h.respondError(ctx, w, http.StatusInternalServerError, "internal server error")
	}
}

// Обработка ошибок, связанных с некорректными запросами
func (h *baseHandler) handleBadRequest(ctx context.Context, w http.ResponseWriter, err error, msg string, args ...any) {
	args = append(args, slog.String("error", err.Error()))

	h.logger.WarnContext(ctx, msg, args...)
	h.respondError(ctx, w, http.StatusBadRequest, msg)
}

// Обработка ошибок валидации
func (h *baseHandler) handleValidationError(ctx context.Context, w http.ResponseWriter, err error, msg string, args ...any) {
	args = append(args, slog.String("error", err.Error()))

	h.logger.WarnContext(ctx, msg, args...)
	h.respondError(ctx, w, http.StatusUnprocessableEntity, msg)
}
