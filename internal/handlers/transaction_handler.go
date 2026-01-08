package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stewicca/angagrar-backend/internal/services"
)

type TransactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

type CreateTransactionRequest struct {
	BudgetID    *uint     `json:"budget_id"`
	Type        string    `json:"type" binding:"required,oneof=income expense"`
	Category    string    `json:"category" binding:"required"`
	Amount      float64   `json:"amount" binding:"required,gt=0"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" binding:"required"`
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.CreateTransaction(
		userID.(uint),
		req.BudgetID,
		req.Type,
		req.Category,
		req.Amount,
		req.Description,
		req.Date,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Transaction created successfully",
		"transaction": transaction,
	})
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	transactions, err := h.transactionService.GetUserTransactions(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
	})
}
