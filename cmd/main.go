package main

import (
	"chat-api/internal/handlers"
	"chat-api/internal/repository"
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=chatdb port=5432 sslmode=disable"
	}

	// 1. Подключение к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 2. Инициализация репозитория и хендлеров
	repo := repository.NewRepository(db)
	h := handlers.NewHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /chats/", h.CreateChat)                  // CreateChat - POST /chats/
	mux.HandleFunc("POST /chats/{id}/messages/", h.CreateMessage) // CreateMessage - POST /chats/{id}/messages/
	mux.HandleFunc("GET /chats/{id}", h.GetChat)                  // GetChat - GET /chats/{id}
	mux.HandleFunc("DELETE /chats/{id}", h.DeleteChat)            // DeleteChat - DELETE /chats/{id}

	loggedMux := LoggingMiddleware(mux)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Server starting on %s\n", port)
	if err := http.ListenAndServe(port, loggedMux); err != nil {
		log.Fatal(err)
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
