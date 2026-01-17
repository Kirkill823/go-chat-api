package models

import (
	"time"
)

type Chat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Messages []Message `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"messages"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ChatID    uint      `gorm:"not null" json:"chat_id"`
	Text      string    `gorm:"type:text;size:5000;not null" json:"text"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
