package service

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"zhanbu/internal/model"
	apperrors "zhanbu/pkg/errors"

	"github.com/rs/zerolog/log"
)

//go:embed prompts/*.txt
var chatPromptTemplates embed.FS

const (
	// MaxHistoryMessages is the maximum number of history messages to include in context.
	MaxHistoryMessages = 20
	// MaxSessionsPerUser is the maximum number of sessions per user.
	MaxSessionsPerUser = 50
	// MaxMessagesPerSession is the maximum number of messages per session.
	MaxMessagesPerSession = 100
)

// ChatRepository defines the interface for chat data access.
type ChatRepository interface {
	CreateSession(session *model.ChatSession) error
	FindSessionByID(id uint) (*model.ChatSession, error)
	FindSessionByUserAndID(userID uint, sessionID uint) (*model.ChatSession, error)
	FindSessionByUserAndRecord(userID uint, recordID uint) (*model.ChatSession, error)
	ListSessionsByUser(userID uint, page, size int) ([]model.ChatSession, int64, error)
	UpdateSessionTitle(id uint, title string) error
	UpdateSessionTimestamp(id uint) error
	DeleteSession(id uint) error
	CreateMessage(message *model.ChatMessage) error
	GetRecentMessages(sessionID uint, limit int) ([]model.ChatMessage, error)
	CountMessagesBySession(sessionID uint) (int64, error)
}

// ChatStartOptions describes a new chat-mode session that is not tied to an
// existing divination record yet.
type ChatStartOptions struct {
	Type     string
	Question string
}

// ChatDivinationStarter creates a real divination record for chat mode.
type ChatDivinationStarter interface {
	Start(userID uint, divinationType string, question string) (*model.DivinationRecord, error)
}

// ChatService handles chat business logic.
type ChatService struct {
	chatRepo         ChatRepository
	recordReader     DivinationRecordReader
	aiProvider       AIProvider
	starter          ChatDivinationStarter
	promptTmpl       string
	liuyaoPromptTmpl string
}

// NewChatService creates a new ChatService.
func NewChatService(chatRepo ChatRepository, recordReader DivinationRecordReader, aiProvider AIProvider) *ChatService {
	svc := &ChatService{
		chatRepo:     chatRepo,
		recordReader: recordReader,
		aiProvider:   aiProvider,
	}

	// Load chat prompt template
	data, err := chatPromptTemplates.ReadFile("prompts/chat_system_prompt.txt")
	if err != nil {
		log.Error().Err(err).Msg("failed to load chat prompt template")
	} else {
		svc.promptTmpl = string(data)
	}

	// Load liuyao v2 prompt template
	liuyaoData, err := chatPromptTemplates.ReadFile("prompts/liuyao_v2_prompt.txt")
	if err != nil {
		log.Error().Err(err).Msg("failed to load liuyao v2 prompt template")
	} else {
		svc.liuyaoPromptTmpl = string(liuyaoData)
	}

	return svc
}

// SetDivinationStarter configures real divination generation for chat mode.
func (s *ChatService) SetDivinationStarter(starter ChatDivinationStarter) {
	s.starter = starter
}

// CreateSession creates a new chat session for a divination record.
func (s *ChatService) CreateSession(userID uint, recordID uint) (*model.ChatSession, error) {
	if s.aiProvider == nil {
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable,
			"AI service is not configured. Set ZHANBU_AI_API_KEY to enable chat.")
	}

	// Verify the record belongs to the user
	record, err := s.recordReader.FindByUserIDAndID(userID, recordID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "divination record not found")
	}

	// Check if session already exists for this record
	existing, _ := s.chatRepo.FindSessionByUserAndRecord(userID, recordID)
	if existing != nil {
		// Load messages and return existing session
		session, err := s.chatRepo.FindSessionByID(existing.ID)
		if err != nil {
			return nil, err
		}
		return session, nil
	}

	// Generate title from question
	title := record.Question
	if title == "" {
		title = fmt.Sprintf("%s占卜", getTypeName(record.Type))
	}
	if len([]rune(title)) > 30 {
		title = string([]rune(title)[:30]) + "..."
	}

	// Create new session
	session := &model.ChatSession{
		UserID:   userID,
		RecordID: recordID,
		Title:    title,
	}
	if err := s.chatRepo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}

	// If there's an existing AI reading, save it as the first assistant message
	if record.AIReading != "" {
		msg := &model.ChatMessage{
			SessionID: session.ID,
			Role:      "assistant",
			Content:   record.AIReading,
		}
		if err := s.chatRepo.CreateMessage(msg); err != nil {
			log.Error().Err(err).Msg("failed to save initial AI reading as message")
		}
	}

	// Reload session with messages
	return s.chatRepo.FindSessionByID(session.ID)
}

// CreateModeSession creates a chat-mode session from a selected divination type
// and the user's first question.
func (s *ChatService) CreateModeSession(userID uint, opts ChatStartOptions) (*model.ChatSession, error) {
	if s.aiProvider == nil {
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable,
			"AI service is not configured. Set ZHANBU_AI_API_KEY to enable chat.")
	}

	divinationType := normalizeChatDivinationType(opts.Type)
	if divinationType == "" {
		return nil, apperrors.New(apperrors.ErrBadRequest, "invalid divination type")
	}
	question := strings.TrimSpace(opts.Question)
	if question == "" {
		return nil, apperrors.New(apperrors.ErrBadRequest, "question is required")
	}

	var record *model.DivinationRecord
	var err error
	if s.starter != nil {
		record, err = s.starter.Start(userID, divinationType, question)
	} else {
		record = &model.DivinationRecord{
			UserID:   userID,
			Type:     divinationType,
			Question: question,
			Result:   buildChatModeInitialResult(divinationType, question),
		}
		err = s.recordReader.Create(record)
	}
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		return nil, fmt.Errorf("failed to create divination record: %w", err)
	}

	title := question
	if len([]rune(title)) > 30 {
		title = string([]rune(title)[:30]) + "..."
	}
	session := &model.ChatSession{
		UserID:   userID,
		RecordID: record.ID,
		Title:    title,
	}
	if err := s.chatRepo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}

	userMsg := &model.ChatMessage{
		SessionID: session.ID,
		Role:      "user",
		Content:   question,
	}
	if err := s.chatRepo.CreateMessage(userMsg); err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	go s.chatRepo.UpdateSessionTimestamp(session.ID)

	return s.chatRepo.FindSessionByID(session.ID)
}

// GetSession gets a chat session with messages.
func (s *ChatService) GetSession(userID uint, sessionID uint) (*model.ChatSession, error) {
	session, err := s.chatRepo.FindSessionByUserAndID(userID, sessionID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "chat session not found")
	}
	return session, nil
}

// ListSessions lists chat sessions for a user.
func (s *ChatService) ListSessions(userID uint, page, size int) ([]model.ChatSession, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}
	return s.chatRepo.ListSessionsByUser(userID, page, size)
}

// DeleteSession deletes a chat session.
func (s *ChatService) DeleteSession(userID uint, sessionID uint) error {
	// Verify ownership
	_, err := s.chatRepo.FindSessionByUserAndID(userID, sessionID)
	if err != nil {
		return apperrors.New(apperrors.ErrNotFound, "chat session not found")
	}
	return s.chatRepo.DeleteSession(sessionID)
}

// SendMessage sends a message and returns a streaming response channel.
func (s *ChatService) SendMessage(userID uint, sessionID uint, content string) (<-chan string, error) {
	if s.aiProvider == nil {
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable,
			"AI service is not configured. Set ZHANBU_AI_API_KEY to enable chat.")
	}

	// Verify session ownership
	session, err := s.chatRepo.FindSessionByUserAndID(userID, sessionID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "chat session not found")
	}

	// Get the divination record
	record, err := s.recordReader.FindByUserIDAndID(userID, session.RecordID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "divination record not found")
	}

	// Check message count limit
	msgCount, _ := s.chatRepo.CountMessagesBySession(sessionID)
	if msgCount >= MaxMessagesPerSession {
		return nil, apperrors.New(apperrors.ErrBadRequest,
			"message limit reached for this session")
	}

	// Save user message
	userMsg := &model.ChatMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   content,
	}
	if err := s.chatRepo.CreateMessage(userMsg); err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Load conversation history
	history, err := s.chatRepo.GetRecentMessages(sessionID, MaxHistoryMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to load conversation history: %w", err)
	}

	// Build messages for AI
	messages := s.buildMessages(record, history)

	// Call AI provider
	ch, err := s.aiProvider.ChatCompletion(messages)
	if err != nil {
		return nil, err
	}

	clientCh := s.streamAndSaveResponse(ch, sessionID, session.RecordID, false)

	// Update session timestamp
	go s.chatRepo.UpdateSessionTimestamp(sessionID)

	return clientCh, nil
}

// StreamInitialReading streams the first AI interpretation for a newly created
// chat-mode session without saving a duplicate user message.
func (s *ChatService) StreamInitialReading(userID uint, sessionID uint) (<-chan string, error) {
	if s.aiProvider == nil {
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable,
			"AI service is not configured. Set ZHANBU_AI_API_KEY to enable chat.")
	}

	session, err := s.chatRepo.FindSessionByUserAndID(userID, sessionID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "chat session not found")
	}

	record, err := s.recordReader.FindByUserIDAndID(userID, session.RecordID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "divination record not found")
	}

	history, err := s.chatRepo.GetRecentMessages(sessionID, MaxHistoryMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to load conversation history: %w", err)
	}
	if len(history) == 0 {
		return nil, apperrors.New(apperrors.ErrBadRequest, "initial question is required")
	}

	ch, err := s.aiProvider.ChatCompletion(s.buildMessages(record, history))
	if err != nil {
		return nil, err
	}

	clientCh := s.streamAndSaveResponse(ch, sessionID, session.RecordID, true)
	go s.chatRepo.UpdateSessionTimestamp(sessionID)

	return clientCh, nil
}

// buildMessages builds the messages array for the AI API.
func (s *ChatService) buildMessages(record *model.DivinationRecord, history []model.ChatMessage) []map[string]string {
	messages := make([]map[string]string, 0, len(history)+1)

	// Build system prompt
	systemPrompt := s.buildSystemPrompt(record)
	messages = append(messages, map[string]string{
		"role":    "system",
		"content": systemPrompt,
	})

	// Add conversation history
	for _, msg := range history {
		messages = append(messages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	return messages
}

// buildSystemPrompt builds the system prompt from the template.
func (s *ChatService) buildSystemPrompt(record *model.DivinationRecord) string {
	// For liuyao_v2, use special prompt with book evidence
	if record.Type == "liuyao_v2" && s.liuyaoPromptTmpl != "" {
		return s.buildLiuyaoPrompt(record)
	}

	// For other types, use general prompt
	if s.promptTmpl == "" {
		// Fallback if template not loaded
		return fmt.Sprintf("你是一位经验丰富的占卜师。用户进行了%s占卜，问题是：%s\n\n占卜结果：\n%s\n\n请基于这个占卜结果为用户提供解读和答疑。",
			record.Type, record.Question, record.Result)
	}

	// Simple template replacement
	prompt := s.promptTmpl
	prompt = strings.ReplaceAll(prompt, "{{.TypeName}}", getTypeName(record.Type))
	prompt = strings.ReplaceAll(prompt, "{{.Question}}", record.Question)
	prompt = strings.ReplaceAll(prompt, "{{.Result}}", record.Result)

	return prompt
}

// buildLiuyaoPrompt builds the prompt for liuyao v2 with book evidence.
func (s *ChatService) buildLiuyaoPrompt(record *model.DivinationRecord) string {
	prompt := s.liuyaoPromptTmpl

	// Parse the result to extract liuyao data
	var reading struct {
		Method       string                       `json:"method"`
		BenGua       *model.TakashimaHexagram     `json:"ben_gua"`
		BianGua      *model.TakashimaHexagram     `json:"bian_gua"`
		MutableLines []int                        `json:"mutable_lines"`
		BookEvidence *model.TakashimaBookEvidence `json:"book_evidence"`
	}
	if err := json.Unmarshal([]byte(record.Result), &reading); err != nil {
		log.Error().Err(err).Msg("failed to parse liuyao result for chat prompt")
		// Fallback to general prompt
		return fmt.Sprintf("你是一位精通《高岛易断》的卦师。用户进行了六爻占卜，问题是：%s\n\n占卜结果：\n%s\n\n请基于这个占卜结果为用户提供解读和答疑。",
			record.Question, record.Result)
	}

	// Fill template
	prompt = strings.ReplaceAll(prompt, "{{.Question}}", record.Question)
	prompt = strings.ReplaceAll(prompt, "{{.Method}}", reading.Method)
	prompt = strings.ReplaceAll(prompt, "{{.BenGua}}", formatPromptHexagram(reading.BenGua))
	prompt = strings.ReplaceAll(prompt, "{{.BianGua}}", formatPromptHexagram(reading.BianGua))
	prompt = strings.ReplaceAll(prompt, "{{.MovingLines}}", formatPromptMovingLines(reading.BenGua, reading.MutableLines))

	if reading.BookEvidence != nil {
		prompt = strings.ReplaceAll(prompt, "{{.BookEvidence}}", formatEvidenceSnippets(reading.BookEvidence.Snippets))
		prompt = strings.ReplaceAll(prompt, "{{.MethodRules}}", formatEvidenceSnippets(reading.BookEvidence.MethodRules))
	} else {
		prompt = strings.ReplaceAll(prompt, "{{.BookEvidence}}", "无")
		prompt = strings.ReplaceAll(prompt, "{{.MethodRules}}", "无")
	}

	return prompt
}

// collectAndSaveResponse collects the full AI response and saves it.
func (s *ChatService) collectAndSaveResponse(ch <-chan string, sessionID uint, recordID uint, saveAsInitialReading bool) {
	var fullResponse strings.Builder
	for chunk := range ch {
		fullResponse.WriteString(chunk)
	}

	s.saveAssistantResponse(sessionID, recordID, fullResponse.String(), saveAsInitialReading)
}

func (s *ChatService) streamAndSaveResponse(ch <-chan string, sessionID uint, recordID uint, saveAsInitialReading bool) <-chan string {
	clientCh := make(chan string, 100)
	go func() {
		defer close(clientCh)

		var fullResponse strings.Builder
		for chunk := range ch {
			fullResponse.WriteString(chunk)
			clientCh <- chunk
		}

		s.saveAssistantResponse(sessionID, recordID, fullResponse.String(), saveAsInitialReading)
	}()
	return clientCh
}

func (s *ChatService) saveAssistantResponse(sessionID uint, recordID uint, response string, saveAsInitialReading bool) {
	if response == "" {
		return
	}

	assistantMsg := &model.ChatMessage{
		SessionID: sessionID,
		Role:      "assistant",
		Content:   response,
	}
	if err := s.chatRepo.CreateMessage(assistantMsg); err != nil {
		log.Error().Err(err).Msg("failed to save assistant message")
	}

	if saveAsInitialReading {
		if err := s.recordReader.UpdateAIReading(recordID, response); err != nil {
			log.Error().Err(err).Uint("record_id", recordID).Msg("failed to save chat response as AI reading")
		}
	}
}

// getTypeName returns the Chinese name for a divination type.
func getTypeName(typeName string) string {
	switch typeName {
	case "tarot":
		return "塔罗牌"
	case "liuyao", "liuyao_v2":
		return "六爻"
	case "bazi":
		return "八字"
	case "horoscope":
		return "星座"
	default:
		return typeName
	}
}

func normalizeChatDivinationType(typeName string) string {
	switch typeName {
	case "tarot", "bazi", "horoscope":
		return typeName
	case "liuyao", "liuyao_v2":
		return "liuyao_v2"
	default:
		return ""
	}
}

func buildChatModeInitialResult(typeName string, question string) string {
	type payload struct {
		Mode     string `json:"mode"`
		Type     string `json:"type"`
		Question string `json:"question"`
	}
	data, _ := json.Marshal(payload{
		Mode:     "chat",
		Type:     typeName,
		Question: question,
	})
	return string(data)
}

// MockChatService is a mock chat service for testing.
type MockChatService struct {
	Session *model.ChatSession
	Err     error
}
