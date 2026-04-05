package handler

import (
	"encoding/json"
	"net/http"

	"task-5/internal/log"
)

type baseHandler struct {
	log log.Logger
}

func NewBaseHandler(l log.Logger) baseHandler {
	return baseHandler{log: l}
}

func (h *baseHandler) respondJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if code != http.StatusNoContent && payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			h.log.Error("failed to encode response", "error", err.Error())
		}
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *baseHandler) respondError(w http.ResponseWriter, code int, message string) {
	h.respondJSON(w, code, errorResponse{message})
}
