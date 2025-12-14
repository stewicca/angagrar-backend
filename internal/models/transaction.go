package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	BudgetID    *uint          `gorm:"index" json:"budget_id,omitempty"`
	Type        string         `gorm:"not null" json:"type"` // income, expense
	Category    string         `gorm:"not null" json:"category"`
	Amount      float64        `gorm:"not null" json:"amount"`
	Description string         `json:"description"`
	Date        time.Time      `gorm:"not null;index" json:"date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
