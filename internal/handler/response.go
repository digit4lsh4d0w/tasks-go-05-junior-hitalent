package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
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
