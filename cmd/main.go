package main

import (
	"fmt"
	"log"
	"net/http"

	"task-5/internal/config"
	"task-5/internal/database"
	"task-5/internal/handler"
	"task-5/internal/logger/slog"
	"task-5/internal/repository/gorm"
	"task-5/internal/service"
)

func main() {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		panic(err)
	}

	// TODO: Remove in production
	fmt.Printf("%#v\n", cfg)

	logger, err := slog.New(&cfg.LogConfig)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	logger.Info("Starting...")

	logger.Debug("Initializing database...")
	db, err := database.NewDatabase(cfg.DBConfig)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		panic(err)
	}

	chatRepo := gorm.NewChatRepositoryGorm(db)
	chatService := service.NewChatService(chatRepo)
	chatHandler := handler.NewChatHandler(chatService)

	http.HandleFunc("POST /chat/", chatHandler.CreateChat)
	http.HandleFunc("POST /chat/{chat_id}/message/", chatHandler.SendMsg)
	http.HandleFunc("GET /chat/{chat_id}/", chatHandler.GetMsgs)
	http.HandleFunc("DELETE /chat/{chat_id}/", chatHandler.DeleteChat)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
