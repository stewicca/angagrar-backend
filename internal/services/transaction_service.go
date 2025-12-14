package services

import (
	"errors"
	"time"

	"github.com/stewicca/angagrar-backend/internal/models"
	"github.com/stewicca/angagrar-backend/internal/repositories"
)

type TransactionService interface {
	CreateTransaction(userID uint, budgetID *uint, transactionType, category string, amount float64, description string, date time.Time) (*models.Transaction, error)
	GetUserTransactions(userID uint) ([]models.Transaction, error)
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository) TransactionService {
	return &transactionService{transactionRepo: transactionRepo}
}

func (s *transactionService) CreateTransaction(userID uint, budgetID *uint, transactionType, category string, amount float64, description string, date time.Time) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	if transactionType != "income" && transactionType != "expense" {
		return nil, errors.New("type must be 'income' or 'expense'")
	}

	transaction := &models.Transaction{
		UserID:      userID,
		BudgetID:    budgetID,
		Type:        transactionType,
		Category:    category,
		Amount:      amount,
		Description: description,
		Date:        date,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *transactionService) GetUserTransactions(userID uint) ([]models.Transaction, error) {
	return s.transactionRepo.FindByUserID(userID)
}
