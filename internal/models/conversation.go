package models

import (
	"time"

	"gorm.io/gorm"
)

type Conversation struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"not null;index" json:"user_id"`
	SessionID       string         `gorm:"uniqueIndex;not null" json:"session_id"` // UUID for frontend reference
	BudgetGenerated bool           `gorm:"default:false" json:"budget_generated"`  // Flag if budget has been generated
	CompletedAt     *time.Time     `json:"completed_at,omitempty"`                 // Null if not completed
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Messages []Message `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}
