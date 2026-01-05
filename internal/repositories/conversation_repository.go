package repositories

import (
	"github.com/stewicca/angagrar-backend/internal/models"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	Create(conversation *models.Conversation) error
	FindByID(id uint) (*models.Conversation, error)
	FindBySessionID(sessionID string) (*models.Conversation, error)
	FindByUserID(userID uint) ([]models.Conversation, error)
	FindActiveByUserID(userID uint) (*models.Conversation, error)
	Update(conversation *models.Conversation) error
	Delete(id uint) error
}

type conversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) Create(conversation *models.Conversation) error {
	return r.db.Create(conversation).Error
}

func (r *conversationRepository) FindByID(id uint) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.db.Preload("Messages").First(&conversation, id).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *conversationRepository) FindBySessionID(sessionID string) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.db.Where("session_id = ?", sessionID).Preload("Messages").First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *conversationRepository) FindByUserID(userID uint) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&conversations).Error
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *conversationRepository) FindActiveByUserID(userID uint) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.db.Where("user_id = ? AND completed_at IS NULL", userID).
		Order("created_at DESC").
		First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *conversationRepository) Update(conversation *models.Conversation) error {
	return r.db.Save(conversation).Error
}

func (r *conversationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Conversation{}, id).Error
}
