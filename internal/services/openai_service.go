package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/stewicca/angagrar-backend/config"
	"github.com/stewicca/angagrar-backend/internal/models"
)

type OpenAIService interface {
	GenerateResponse(systemPrompt string, messages []models.Message) (string, error)
	GenerateResponseWithRetry(systemPrompt string, messages []models.Message, maxRetries int) (string, error)
}

type openAIService struct {
	client      *openai.Client
	model       string
	maxTokens   int
	temperature float32
}

func NewOpenAIService(cfg *config.Config) OpenAIService {
	client := openai.NewClient(cfg.OpenAIAPIKey)

	return &openAIService{
		client:      client,
		model:       cfg.OpenAIModel,
		maxTokens:   cfg.OpenAIMaxTokens,
		temperature: cfg.OpenAITemp,
	}
}

// GenerateResponse calls OpenAI API to generate a response
func (s *openAIService) GenerateResponse(systemPrompt string, messages []models.Message) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert messages to OpenAI format
	chatMessages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
	}

	// Add conversation history
	for _, msg := range messages {
		role := openai.ChatMessageRoleUser
		if msg.Role == models.RoleAssistant {
			role = openai.ChatMessageRoleAssistant
		}

		chatMessages = append(chatMessages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Create chat completion request
	req := openai.ChatCompletionRequest{
		Model:       s.model,
		Messages:    chatMessages,
		MaxTokens:   s.maxTokens,
		Temperature: s.temperature,
	}

	// Call OpenAI API
	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateResponseWithRetry attempts to generate response with exponential backoff retry
func (s *openAIService) GenerateResponseWithRetry(systemPrompt string, messages []models.Message, maxRetries int) (string, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		response, err := s.GenerateResponse(systemPrompt, messages)
		if err == nil {
			return response, nil
		}

		lastErr = err

		// Exponential backoff: 1s, 2s, 4s
		if attempt < maxRetries-1 {
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			time.Sleep(backoff)
		}
	}

	return "", fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
