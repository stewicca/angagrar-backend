package repositories

import (
	"github.com/stewicca/angagrar-backend/internal/models"
	"gorm.io/gorm"
)

type BudgetRepository interface {
	Create(budget *models.Budget) error
	CreateBatch(budgets []models.Budget) error
	FindByID(id uint) (*models.Budget, error)
	FindByUserID(userID uint) ([]models.Budget, error)
	Update(budget *models.Budget) error
	Delete(id uint) error
}

type budgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) Create(budget *models.Budget) error {
	return r.db.Create(budget).Error
}

func (r *budgetRepository) CreateBatch(budgets []models.Budget) error {
	return r.db.Create(&budgets).Error
}

func (r *budgetRepository) FindByID(id uint) (*models.Budget, error) {
	var budget models.Budget
	err := r.db.First(&budget, id).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *budgetRepository) FindByUserID(userID uint) ([]models.Budget, error) {
	var budgets []models.Budget
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&budgets).Error
	if err != nil {
		return nil, err
	}
	return budgets, nil
}

func (r *budgetRepository) Update(budget *models.Budget) error {
	return r.db.Save(budget).Error
}

func (r *budgetRepository) Delete(id uint) error {
	return r.db.Delete(&models.Budget{}, id).Error
}
