package repository

import (
	"chat-api/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateChat(chat *models.Chat) error {
	return r.db.Create(chat).Error
}

func (r *Repository) CreateMessage(msg *models.Message) error {
	// Сначала проверим, существует ли чат с таким ID
	var count int64
	if err := r.db.Model(&models.Chat{}).Where("id = ?", msg.ChatID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.Create(msg).Error
}

func (r *Repository) GetChatWithMessages(chatID int, limit int) (*models.Chat, error) {
	var chat models.Chat

	// 1. Ищем чат
	if err := r.db.First(&chat, chatID).Error; err != nil {
		return nil, err
	}
	// 2. Ищем сообщения этого чата с лимитом
	err := r.db.Model(&models.Message{}).
		Where("chat_id = ?", chatID).
		Order("created_at desc").
		Limit(limit).
		Find(&chat.Messages).Error

	return &chat, err
}

func (r *Repository) DeleteChat(chatID int) error {
	result := r.db.Delete(&models.Chat{}, chatID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
