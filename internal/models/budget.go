package models

import (
	"time"

	"gorm.io/gorm"
)

type Budget struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	Category     string         `gorm:"not null" json:"category"`
	Amount       float64        `gorm:"not null" json:"amount"`
	Period       string         `gorm:"not null" json:"period"` // monthly, yearly
	StartDate    time.Time      `gorm:"not null" json:"start_date"`
	EndDate      time.Time      `gorm:"not null" json:"end_date"`
	Description  string         `json:"description"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Transactions []Transaction  `gorm:"foreignKey:BudgetID" json:"transactions,omitempty"`
}
