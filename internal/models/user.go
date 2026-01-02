package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a guest user
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	GuestID   string         `gorm:"uniqueIndex;not null" json:"guest_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Budgets   []Budget       `gorm:"foreignKey:UserID" json:"budgets,omitempty"`
}
