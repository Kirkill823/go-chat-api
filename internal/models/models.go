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

type User struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	FirstName  string  `gorm:"not null" json:"first_name"`
	LastName   string  `gorm:"not null" json:"last_name"`
	SecondName *string `json:"second_name,omitempty"`
	Password   string  `gorm:"not null" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
