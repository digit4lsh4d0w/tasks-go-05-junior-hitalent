package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"task-5/internal/log"
	"task-5/internal/model"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

const (
	limitDefault = 20
	limitMax     = 100
)

type ChatService interface {
	GetChat(id uint, limit int) (*model.Chat, error)
	CreateChat(chat *model.Chat) error
	DeleteChat(id uint) error
}

type chatHandler struct {
	baseHandler
	service   ChatService
	validator *validator.Validate
}

func NewChatHandler(s ChatService, v *validator.Validate, l log.Logger) *chatHandler {
	return &chatHandler{
		baseHandler: NewBaseHandler(l),
		service:     s,
		validator:   v,
	}
}

type CreateChatRequest struct {
	Title string `json:"title" validate:"required,min=1,max=200"`
}

func (h *chatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	chat := &model.Chat{
		Title: req.Title,
	}

	if err := h.service.CreateChat(chat); err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to create chat")
		return
	}

	h.respondJSON(w, http.StatusCreated, chat)
}

func (h *chatHandler) SendMsg(w http.ResponseWriter, r *http.Request) {}

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

func (h *chatHandler) GetMsgs(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chat_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid chat id")
		return
	}

	limit := parseLimit(r)

	chat, err := h.service.GetChat(uint(id), limit)
	if err != nil {
		// TODO: Исправить протечку абстракции
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.respondError(w, http.StatusNotFound, "not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.respondJSON(w, http.StatusOK, chat)
}

func (h *chatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {}
