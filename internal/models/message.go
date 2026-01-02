package models

import (
	"time"

	"gorm.io/gorm"
)

// MessageRole represents who sent the message
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

type Message struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ConversationID uint           `gorm:"not null;index" json:"conversation_id"`
	Role           MessageRole    `gorm:"not null" json:"role"` // "user" or "assistant"
	Content        string         `gorm:"type:text;not null" json:"content"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Conversation Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
}
