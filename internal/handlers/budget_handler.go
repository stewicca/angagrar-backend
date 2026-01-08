package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stewicca/angagrar-backend/internal/repositories"
	"github.com/stewicca/angagrar-backend/pkg/utils"
)

type BudgetHandler struct {
	budgetRepo repositories.BudgetRepository
}

func NewBudgetHandler(budgetRepo repositories.BudgetRepository) *BudgetHandler {
	return &BudgetHandler{
		budgetRepo: budgetRepo,
	}
}

// GetUserBudgets handles GET /api/v1/budgets
func (h *BudgetHandler) GetUserBudgets(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	budgets, err := h.budgetRepo.FindByUserID(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve budgets", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Budgets retrieved", gin.H{
		"budgets": budgets,
	})
}

// UpdateBudget handles PATCH /api/v1/budgets/:id
func (h *BudgetHandler) UpdateBudget(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	budgetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid budget ID")
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Amount is required")
		return
	}

	// Find budget
	budget, err := h.budgetRepo.FindByID(uint(budgetID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Budget not found", err)
		return
	}

	// Verify ownership
	if budget.UserID != userID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "You don't have permission to update this budget", nil)
		return
	}

	// Validate amount
	if req.Amount <= 0 {
		utils.ValidationErrorResponse(c, "Amount must be greater than 0")
		return
	}

	// Update amount
	budget.Amount = req.Amount
	if err := h.budgetRepo.Update(budget); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update budget", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Budget updated", gin.H{
		"budget": budget,
	})
}
