package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stewicca/angagrar-backend/internal/services"
	"github.com/stewicca/angagrar-backend/pkg/utils"
)

type ConversationHandler struct {
	conversationService services.ConversationService
}

func NewConversationHandler(conversationService services.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
	}
}

// StartConversation handles POST /api/v1/conversations/start
func (h *ConversationHandler) StartConversation(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	conversation, greetingMsg, err := h.conversationService.StartConversation(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to start conversation", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Conversation started", gin.H{
		"session_id": conversation.SessionID,
		"message":    greetingMsg,
	})
}

// SendMessage handles POST /api/v1/conversations/:sessionId/messages
func (h *ConversationHandler) SendMessage(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Message is required")
		return
	}

	response, isCompleted, budgets, err := h.conversationService.ProcessMessage(sessionID, req.Message)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to process message", err)
		return
	}

	responseData := gin.H{
		"assistant_message": response,
		"completed":         isCompleted,
	}

	// Include budgets if generated
	if budgets != nil && len(budgets) > 0 {
		responseData["budgets"] = budgets
		responseData["budget_generated"] = true
	}

	utils.SuccessResponse(c, http.StatusOK, "Message processed", responseData)
}

// GetConversationHistory handles GET /api/v1/conversations/:sessionId/history
func (h *ConversationHandler) GetConversationHistory(c *gin.Context) {
	sessionID := c.Param("sessionId")

	messages, err := h.conversationService.GetConversationHistory(sessionID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Conversation not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "History retrieved", gin.H{
		"messages": messages,
	})
}

// ResetConversation handles POST /api/v1/conversations/:sessionId/reset
func (h *ConversationHandler) ResetConversation(c *gin.Context) {
	sessionID := c.Param("sessionId")

	conversation, greetingMsg, err := h.conversationService.ResetConversation(sessionID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to reset conversation", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Conversation reset", gin.H{
		"message":        "Conversation reset. Starting new interview.",
		"new_session_id": conversation.SessionID,
		"greeting":       greetingMsg,
	})
}
