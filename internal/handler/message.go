package handler

type MessageService interface{}

type CreateMessageRequest struct {
	ChatID uint
	Text   string `validate:"required,min=1,max=5000"`
}
