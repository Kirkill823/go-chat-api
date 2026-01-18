package handlers

import (
	"chat-api/internal/models"
	"chat-api/internal/repository"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// CreateChat - POST /chats/
func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if len(req.Title) == 0 || len(req.Title) > 200 {
		http.Error(w, "Title must be between 1 and 200 chars", http.StatusBadRequest)
		return
	}

	chat := models.Chat{Title: req.Title}
	if err := h.repo.CreateChat(&chat); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// CreateMessage - POST /chats/{id}/messages/
func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	chatID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Chat ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидация
	if len(req.Text) == 0 || len(req.Text) > 5000 {
		http.Error(w, "Text must be between 1 and 5000 chars", http.StatusBadRequest)
		return
	}

	msg := models.Message{
		ChatID: uint(chatID),
		Text:   req.Text,
	}

	if err := h.repo.CreateMessage(&msg); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// GetChat - GET /chats/{id}
func (h *Handler) GetChat(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	chatID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Chat ID", http.StatusBadRequest)
		return
	}
	// Лимит сообщений (по умолчанию 20, максимум 100)
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		parsedLimit, _ := strconv.Atoi(l)
		if parsedLimit > 0 {
			limit = parsedLimit
		}
		if parsedLimit > 100 {
			limit = 100
		}
	}

	chat, err := h.repo.GetChatWithMessages(chatID, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// DeleteChat - DELETE /chats/{id}
func (h *Handler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	chatID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Chat ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteChat(chatID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
