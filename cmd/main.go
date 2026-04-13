package main

import (
	"log/slog"
	"net/http"

	"task-5/internal/config"
	"task-5/internal/db"
	"task-5/internal/handler"
	"task-5/internal/logger"
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

	logger, cleanup, err := logger.New(&cfg.LogConfig)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	slog.SetDefault(logger)

	logger.Info(
		"Starting...",
		slog.String("endpoint", cfg.Endpoint),
	)

	db, err := db.NewDatabase(cfg.DBConfig)
	if err != nil {
		logger.Error("Failed to connect to database", slog.String("error", err.Error()))
		panic(err)
	}
	logger.Info("Connected to database")

	validator := validator.New()

	chatRepo := gorm.NewChatRepository(db)
	chatService := service.NewChatService(
		chatRepo,
		logger.With(slog.String("layer", "service")),
	)
	chatHandler := handler.NewChatHandler(
		chatService,
		validator,
		logger.With(slog.String("layer", "transport")),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /chat/", chatHandler.CreateChat)
	mux.HandleFunc("POST /chat/{chat_id}/message/", chatHandler.CreateMessage)
	mux.HandleFunc("GET /chat/{chat_id}/", chatHandler.GetAllMessages)
	mux.HandleFunc("DELETE /chat/{chat_id}/", chatHandler.DeleteChat)

	handler := middleware.Chain(
		mux,
		middleware.TraceID(),
		middleware.Logger(logger.With(slog.String("layer", "middleware"))),
	)

	err = http.ListenAndServe(cfg.Endpoint, handler)
	if err != nil && err != http.ErrServerClosed {
		logger.Error("Server stopped", slog.String("error", err.Error()))
	}
}
