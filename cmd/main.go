package main

import (
	"fmt"
	"net/http"

	"task-5/internal/config"
	"task-5/internal/db"
	"task-5/internal/handler"
	"task-5/internal/log/slog"
	"task-5/internal/middleware"
	"task-5/internal/repository/gorm"
	"task-5/internal/service"

	"github.com/go-playground/validator"
)

func main() {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		panic(err)
	}

	// TODO: Remove in production
	fmt.Printf("%#v\n", cfg)

	log, err := slog.New(&cfg.LogConfig)
	if err != nil {
		panic(err)
	}
	defer log.Close()

	log.Info("Starting...")

	log.Debug("Initializing database...")
	db, err := db.NewDatabase(cfg.DBConfig)
	if err != nil {
		log.Error("Failed to initialize database", "error", err)
		panic(err)
	}

	validator := validator.New()

	chatRepo := gorm.NewChatRepository(db, log)
	msgRepo := gorm.NewMessageRepository(db, log)
	chatService := service.NewChatService(chatRepo, msgRepo)
	chatHandler := handler.NewChatHandler(chatService, validator, log.With("handler", "chat"))

	mux := http.NewServeMux()
	mux.HandleFunc("POST /chat/", chatHandler.CreateChat)
	mux.HandleFunc("POST /chat/{chat_id}/message/", chatHandler.CreateMessage)
	mux.HandleFunc("GET /chat/{chat_id}/", chatHandler.GetAllMessages)
	mux.HandleFunc("DELETE /chat/{chat_id}/", chatHandler.DeleteChat)

	handler := middleware.Chain(mux, middleware.Log(log))

	log.Error("Server stopped", "error", http.ListenAndServe(":3000", handler))
}
