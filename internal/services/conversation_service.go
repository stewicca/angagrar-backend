package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/stewicca/angagrar-backend/internal/models"
	"github.com/stewicca/angagrar-backend/internal/repositories"
)

type ConversationService interface {
	StartConversation(userID uint) (*models.Conversation, string, error)
	ProcessMessage(sessionID string, userMessage string) (string, bool, []models.Budget, error)
	GetConversationHistory(sessionID string) ([]models.Message, error)
	ResetConversation(sessionID string) (*models.Conversation, string, error)
}

type conversationService struct {
	conversationRepo repositories.ConversationRepository
	messageRepo      repositories.MessageRepository
	budgetRepo       repositories.BudgetRepository
	openAIService    OpenAIService
}

func NewConversationService(
	conversationRepo repositories.ConversationRepository,
	messageRepo repositories.MessageRepository,
	budgetRepo repositories.BudgetRepository,
	openAIService OpenAIService,
) ConversationService {
	return &conversationService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		budgetRepo:       budgetRepo,
		openAIService:    openAIService,
	}
}

// StartConversation creates a new conversation and returns greeting message
func (s *conversationService) StartConversation(userID uint) (*models.Conversation, string, error) {
	// Check if user already has an active conversation (for MVP: 1 conversation only)
	existingConv, err := s.conversationRepo.FindActiveByUserID(userID)
	if err == nil && existingConv != nil {
		return nil, "", fmt.Errorf("user already has an active conversation, complete or reset it first")
	}

	// Create new conversation
	sessionID := uuid.New().String()
	conversation := &models.Conversation{
		UserID:          userID,
		SessionID:       sessionID,
		BudgetGenerated: false,
	}

	if err := s.conversationRepo.Create(conversation); err != nil {
		return nil, "", fmt.Errorf("failed to create conversation: %w", err)
	}

	// Generate personalized greeting using LLM
	systemPrompt := getAiraSystemPrompt()
	initialMessages := []models.Message{}

	greetingMsg, err := s.openAIService.GenerateResponseWithRetry(systemPrompt, initialMessages, 3)
	if err != nil {
		// Fallback greeting
		greetingMsg = "hai! 👋 gue aira, siap bantu kamu atur budget yang pas buat lifestyle kamu. cerita aja dulu tentang keuangan kamu, gaji berapa, tinggal dimana, lifestyle gimana?"
	}

	// Save assistant message
	assistantMsg := &models.Message{
		ConversationID: conversation.ID,
		Role:           models.RoleAssistant,
		Content:        greetingMsg,
	}
	if err := s.messageRepo.Create(assistantMsg); err != nil {
		return nil, "", fmt.Errorf("failed to save message: %w", err)
	}

	return conversation, greetingMsg, nil
}

// ProcessMessage handles user input and generates AI response
func (s *conversationService) ProcessMessage(sessionID string, userMessage string) (string, bool, []models.Budget, error) {
	// Find conversation
	conversation, err := s.conversationRepo.FindBySessionID(sessionID)
	if err != nil {
		return "", false, nil, fmt.Errorf("conversation not found: %w", err)
	}

	// Check if conversation is already completed
	if conversation.CompletedAt != nil {
		return "Conversation sudah selesai. Silakan start conversation baru.", true, nil, nil
	}

	// Save user message
	userMsg := &models.Message{
		ConversationID: conversation.ID,
		Role:           models.RoleUser,
		Content:        userMessage,
	}
	if err := s.messageRepo.Create(userMsg); err != nil {
		return "", false, nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Get conversation history for context
	messages, err := s.messageRepo.FindByConversationID(conversation.ID)
	if err != nil {
		return "", false, nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	// Check if user asks to generate budget
	shouldGenerateBudget := s.detectBudgetGenerationIntent(userMessage, messages)

	var budgets []models.Budget
	var aiResponse string

	if shouldGenerateBudget && !conversation.BudgetGenerated {
		// Ask LLM to analyze conversation and generate budget
		budgets, aiResponse, err = s.generateBudgetFromConversation(conversation, messages)
		if err != nil {
			return "maaf, ada error saat generate budget 😅 coba lagi ya!", false, nil, err
		}

		// Mark budget as generated
		conversation.BudgetGenerated = true
		now := time.Now()
		conversation.CompletedAt = &now
		if err := s.conversationRepo.Update(conversation); err != nil {
			return "", false, nil, fmt.Errorf("failed to update conversation: %w", err)
		}
	} else {
		// Continue conversation normally
		systemPrompt := getAiraSystemPrompt()
		aiResponse, err = s.openAIService.GenerateResponseWithRetry(systemPrompt, messages, 3)
		if err != nil {
			aiResponse = "hmm gue lagi error nih 😅 bisa coba lagi?"
		}
	}

	// Save assistant response
	assistantMsg := &models.Message{
		ConversationID: conversation.ID,
		Role:           models.RoleAssistant,
		Content:        aiResponse,
	}
	if err := s.messageRepo.Create(assistantMsg); err != nil {
		return "", false, nil, fmt.Errorf("failed to save assistant message: %w", err)
	}

	isCompleted := conversation.CompletedAt != nil
	return aiResponse, isCompleted, budgets, nil
}

// detectBudgetGenerationIntent checks if user wants to generate budget
func (s *conversationService) detectBudgetGenerationIntent(userMessage string, messages []models.Message) bool {
	lower := strings.ToLower(userMessage)

	// Keywords that trigger budget generation
	keywords := []string{
		"buatin budget",
		"bikinin budget",
		"generate budget",
		"buat budget",
		"siap",
		"oke buatin",
		"lanjut",
		"udah cukup",
	}

	for _, keyword := range keywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}

	// If conversation is long enough (>= 6 messages), assume ready for budget
	if len(messages) >= 6 {
		return true
	}

	return false
}

// generateBudgetFromConversation uses LLM to analyze conversation and generate personalized budget
func (s *conversationService) generateBudgetFromConversation(conversation *models.Conversation, messages []models.Message) ([]models.Budget, string, error) {
	// Create prompt for LLM to analyze conversation and generate budget
	analysisPrompt := getBudgetAnalysisPrompt(messages)

	// Call LLM to get budget recommendation
	llmResponse, err := s.openAIService.GenerateResponseWithRetry(analysisPrompt, []models.Message{}, 3)
	if err != nil {
		return nil, "", fmt.Errorf("LLM analysis failed: %w", err)
	}

	// Parse LLM response to extract budget data
	budgetData, err := s.parseLLMBudgetResponse(llmResponse)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse budget response: %w", err)
	}

	// Create budget records
	budgets := s.createBudgetRecords(conversation.UserID, budgetData)

	// Save budgets to database
	if err := s.budgetRepo.CreateBatch(budgets); err != nil {
		return nil, "", fmt.Errorf("failed to save budgets: %w", err)
	}

	// Generate user-friendly response
	response := formatBudgetResponse(budgets, budgetData)

	return budgets, response, nil
}

// GetConversationHistory retrieves all messages in a conversation
func (s *conversationService) GetConversationHistory(sessionID string) ([]models.Message, error) {
	conversation, err := s.conversationRepo.FindBySessionID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}

	messages, err := s.messageRepo.FindByConversationID(conversation.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}

	return messages, nil
}

// ResetConversation resets an existing conversation
func (s *conversationService) ResetConversation(sessionID string) (*models.Conversation, string, error) {
	conversation, err := s.conversationRepo.FindBySessionID(sessionID)
	if err != nil {
		return nil, "", fmt.Errorf("conversation not found: %w", err)
	}

	// Delete old conversation
	if err := s.conversationRepo.Delete(conversation.ID); err != nil {
		return nil, "", fmt.Errorf("failed to delete conversation: %w", err)
	}

	// Start new conversation
	return s.StartConversation(conversation.UserID)
}

// Helper functions

func getAiraSystemPrompt() string {
	return `Kamu adalah Aira, asisten keuangan virtual yang asik dan helpful banget!

PERSONALITY:
- Tone: Santai, Gen Z Indonesia, casual tapi tetap sopan
- Style: Ramah, encouraging, tidak judgemental
- Language: Bahasa Indonesia casual (gue/kamu, bukan saya/anda)
- Max 2-3 kalimat per response, jangan bertele-tele
- Gunakan emoji secukupnya (max 2 per message)

TUGAS KAMU:
Kamu membantu user membuat budget personal yang cocok untuk mereka. Obrolan kamu natural dan tidak kaku.

YANG PERLU KAMU CARI TAHU (tapi jangan kaku, natural aja):
1. Gaji/income bulanan mereka
2. Lokasi tinggal (kota mana)
3. Lifestyle preference (hemat, moderate, atau santai)
4. Pengeluaran rutin apa aja
5. Goals keuangan (nabung, invest, dll)
6. Kebiasaan spending (sering healing, hobi mahal, dll)

ATURAN:
- Jangan tanya semua sekaligus, ngobrol natural
- Kalau user udah cerita banyak, tawarkan untuk bikinin budget
- Jangan judgmental, supportive aja
- Kalau user bilang "buatin budget" atau sejenisnya, artinya mereka siap

GREETING PERTAMA:
Sambut user dengan ramah dan ajak mereka cerita tentang keuangan mereka secara casual.`
}

func getBudgetAnalysisPrompt(messages []models.Message) string {
	// Convert messages to conversation transcript
	transcript := ""
	for _, msg := range messages {
		role := "User"
		if msg.Role == models.RoleAssistant {
			role = "Aira"
		}
		transcript += fmt.Sprintf("%s: %s\n", role, msg.Content)
	}

	return fmt.Sprintf(`Kamu adalah AI budget analyst. Analisa percakapan berikut dan generate personalized budget.

PERCAKAPAN:
%s

TUGAS KAMU:
1. Extract informasi penting: salary, location, lifestyle, spending habits, goals
2. Pertimbangkan cost of living di lokasi mereka
3. Pertimbangkan lifestyle dan kebiasaan mereka
4. Generate budget allocation yang PERSONAL dan REALISTIC

OUTPUT FORMAT (JSON):
{
  "salary": <angka>,
  "location": "<kota>",
  "analysis": "<penjelasan singkat kenapa budget ini cocok untuk mereka>",
  "categories": [
    {"name": "Kewajiban", "amount": <angka>, "description": "sewa, utilities, cicilan"},
    {"name": "Makan", "amount": <angka>, "description": "makanan sehari-hari"},
    {"name": "Transport", "amount": <angka>, "description": "transportasi"},
    {"name": "Healing", "amount": <angka>, "description": "hiburan, self-care"},
    {"name": "Tabungan", "amount": <angka>, "description": "tabungan & investasi"},
    {"name": "Lain-lain", "amount": <angka>, "description": "pengeluaran lain"}
  ]
}

PENTING:
- Total semua amount HARUS = salary
- Round ke nearest 1000
- Realistic dengan cost of living kota mereka
- Personal based on habits & goals mereka

Return ONLY valid JSON, no explanation.`, transcript)
}

type BudgetData struct {
	Salary     float64 `json:"salary"`
	Location   string  `json:"location"`
	Analysis   string  `json:"analysis"`
	Categories []struct {
		Name        string  `json:"name"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	} `json:"categories"`
}

func (s *conversationService) parseLLMBudgetResponse(llmResponse string) (*BudgetData, error) {
	// Extract JSON from response (LLM might add extra text)
	start := strings.Index(llmResponse, "{")
	end := strings.LastIndex(llmResponse, "}") + 1

	if start == -1 || end == 0 {
		return nil, fmt.Errorf("no JSON found in response")
	}

	jsonStr := llmResponse[start:end]

	var budgetData BudgetData
	if err := json.Unmarshal([]byte(jsonStr), &budgetData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &budgetData, nil
}

func (s *conversationService) createBudgetRecords(userID uint, data *BudgetData) []models.Budget {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	budgets := []models.Budget{}
	for _, cat := range data.Categories {
		budgets = append(budgets, models.Budget{
			UserID:      userID,
			Category:    cat.Name,
			Amount:      cat.Amount,
			Period:      "monthly",
			StartDate:   startDate,
			EndDate:     endDate,
			Description: cat.Description,
		})
	}

	return budgets
}

func formatBudgetResponse(budgets []models.Budget, data *BudgetData) string {
	response := "done! ✨ ini budget recommendation yang gue bikinin buat kamu:\n\n"

	for _, b := range budgets {
		emoji := getCategoryEmoji(b.Category)
		response += fmt.Sprintf("%s %s: Rp %.0f\n", emoji, b.Category, b.Amount)
	}

	response += fmt.Sprintf("\n💡 %s\n\n", data.Analysis)
	response += "kamu bisa adjust sendiri nanti kalau ada yang kurang pas!"

	return response
}

func getCategoryEmoji(category string) string {
	emojis := map[string]string{
		"Kewajiban": "💸",
		"Makan":     "🍜",
		"Transport": "🚗",
		"Healing":   "🎮",
		"Tabungan":  "💰",
		"Lain-lain": "📦",
	}

	if emoji, ok := emojis[category]; ok {
		return emoji
	}
	return "💵"
}
