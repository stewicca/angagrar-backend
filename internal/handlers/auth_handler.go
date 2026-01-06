package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stewicca/angagrar-backend/internal/services"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) CreateGuest(c *gin.Context) {
	user, token, err := h.authService.CreateGuest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create guest user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}
