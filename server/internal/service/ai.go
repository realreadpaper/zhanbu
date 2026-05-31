package service

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"

	"zhanbu/config"
	"zhanbu/internal/model"
	apperrors "zhanbu/pkg/errors"
)

//go:embed prompts/*.txt
var promptTemplates embed.FS

// AIProvider defines the interface for AI service providers.
type AIProvider interface {
	// Interpret generates an AI interpretation for a divination result.
	// It should return a channel that streams text chunks.
	Interpret(divinationType string, result string, question string) (<-chan string, error)
}

// OpenAIProvider implements AIProvider using OpenAI-compatible API.
type OpenAIProvider struct {
	apiKey      string
	baseURL     string
	model       string
	maxTokens   int
	temperature float64
	client      *http.Client
	prompts     map[string]*template.Template
}

// NewOpenAIProvider creates a new OpenAI-compatible AI provider.
func NewOpenAIProvider(cfg *config.AIConfig) (*OpenAIProvider, error) {
	provider := &OpenAIProvider{
		apiKey:      cfg.APIKey,
		baseURL:     cfg.BaseURL,
		model:       cfg.Model,
		maxTokens:   cfg.MaxTokens,
		temperature: cfg.Temperature,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		prompts: make(map[string]*template.Template),
	}

	if err := provider.loadPrompts(); err != nil {
		return nil, fmt.Errorf("failed to load prompt templates: %w", err)
	}

	return provider, nil
}

// loadPrompts loads prompt templates from embedded files.
func (p *OpenAIProvider) loadPrompts() error {
	types := []string{"tarot", "horoscope", "liuyao", "bazi"}
	for _, t := range types {
		filename := fmt.Sprintf("prompts/%s_prompt.txt", t)
		data, err := promptTemplates.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read prompt template %s: %w", filename, err)
		}
		tmpl, err := template.New(t).Parse(string(data))
		if err != nil {
			return fmt.Errorf("failed to parse prompt template %s: %w", filename, err)
		}
		p.prompts[t] = tmpl
	}
	return nil
}

// Interpret streams an AI interpretation for a divination result.
func (p *OpenAIProvider) Interpret(divinationType string, result string, question string) (<-chan string, error) {
	tmpl, ok := p.prompts[divinationType]
	if !ok {
		return nil, apperrors.New(apperrors.ErrBadRequest, fmt.Sprintf("unsupported divination type: %s", divinationType))
	}

	// Build prompt from template
	var promptBuf bytes.Buffer
	err := tmpl.Execute(&promptBuf, map[string]string{
		"Spread":  extractSpread(result),
		"Zodiac":  extractZodiac(result),
		"Period":  extractPeriod(result),
		"Question": question,
		"Result":  result,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute prompt template: %w", err)
	}

	// Make streaming API request
	reqBody := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "你是一位经验丰富的占卜师，精通东方和西方占卜术。请用温暖亲切的语言为用户提供专业的解读。",
			},
			{
				"role":    "user",
				"content": promptBuf.String(),
			},
		},
		"max_tokens":  p.maxTokens,
		"temperature": p.temperature,
		"stream":      true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := strings.TrimRight(p.baseURL, "/") + "/chat/completions"
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, apperrors.NewWithErr(apperrors.ErrAIServiceUnavailable, "AI service request failed", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable,
			fmt.Sprintf("AI service returned status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse streaming response
	ch := make(chan string, 100)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}

			var chunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				ch <- chunk.Choices[0].Delta.Content
			}
		}
	}()

	return ch, nil
}

// extractSpread extracts spread name from tarot result JSON.
func extractSpread(result string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		return ""
	}
	if s, ok := m["spread"].(string); ok {
		return s
	}
	return ""
}

// extractZodiac extracts zodiac name from horoscope result JSON.
func extractZodiac(result string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		return ""
	}
	if z, ok := m["zodiac"].(string); ok {
		return z
	}
	return ""
}

// extractPeriod extracts period from horoscope result JSON.
func extractPeriod(result string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		return ""
	}
	if p, ok := m["period"].(string); ok {
		return p
	}
	return ""
}

// AIResultReader reads divination records for AI interpretation.
type AIResultReader interface {
	FindByUserIDAndID(userID uint, id uint) (*model.DivinationRecord, error)
	UpdateAIReading(id uint, reading string) error
}

// AIService handles AI interpretation business logic.
type AIService struct {
	provider AIProvider
	reader   AIResultReader
}

// NewAIService creates a new AIService.
func NewAIService(provider AIProvider, reader AIResultReader) *AIService {
	return &AIService{
		provider: provider,
		reader:   reader,
	}
}

// Interpret starts streaming AI interpretation for a divination record.
func (s *AIService) Interpret(userID uint, recordID uint, divinationType string, question string) (<-chan string, error) {
	if s.provider == nil {
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable,
			"AI service is not configured. Set ZHANBU_AI_API_KEY to enable AI readings.")
	}

	// Verify the record belongs to the user
	record, err := s.reader.FindByUserIDAndID(userID, recordID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "divination record not found")
	}

	// Use the record's type if not specified
	if divinationType == "" {
		divinationType = record.Type
	}

	// Use the record's question if not specified
	if question == "" {
		question = record.Question
	}

	// Check if AI reading already exists
	if record.AIReading != "" {
		// Return existing reading as a single chunk
		ch := make(chan string, 1)
		ch <- record.AIReading
		close(ch)
		return ch, nil
	}

	// Start streaming interpretation
	ch, err := s.provider.Interpret(divinationType, record.Result, question)
	if err != nil {
		return nil, err
	}

	// Collect the full reading in background for saving
	fullReadingCh := make(chan string, 100)
	go func() {
		var fullReading strings.Builder
		for chunk := range ch {
			fullReading.WriteString(chunk)
			fullReadingCh <- chunk
		}
		close(fullReadingCh)

		// Save the complete reading
		if fullReading.Len() > 0 {
			_ = s.reader.UpdateAIReading(recordID, fullReading.String())
		}
	}()

	return fullReadingCh, nil
}

// InterpretDirect starts streaming AI interpretation with direct result data (no record).
func (s *AIService) InterpretDirect(divinationType string, result string, question string) (<-chan string, error) {
	if s.provider == nil {
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable,
			"AI service is not configured. Set ZHANBU_AI_API_KEY to enable AI readings.")
	}

	return s.provider.Interpret(divinationType, result, question)
}

// MockAIProvider is a mock AI provider for testing.
type MockAIProvider struct {
	Response string
	Err      error
}

// Interpret returns a mock response.
func (m *MockAIProvider) Interpret(divinationType string, result string, question string) (<-chan string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	ch := make(chan string, 100)
	go func() {
		defer close(ch)
		// Send the full response as a single chunk for reliable testing
		if m.Response != "" {
			ch <- m.Response
		}
	}()
	return ch, nil
}
