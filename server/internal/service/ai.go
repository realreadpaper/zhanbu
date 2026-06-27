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

	"github.com/rs/zerolog/log"
)

//go:embed prompts/*.txt
var promptTemplates embed.FS

const (
	aiHTTPTimeout              = 300 * time.Second
	incompleteLengthNotice     = "\n\n【系统提示：AI 输出达到长度上限，内容可能未完整生成，请稍后重新解读，或调高 AI_MAX_TOKENS 后重试。】"
	incompleteConnectionNotice = "\n\n【系统提示：AI 解读连接中断，内容可能未完整生成，请稍后重新解读。】"
)

// AIProvider defines the interface for AI service providers.
type AIProvider interface {
	// Interpret generates an AI interpretation for a divination result.
	// It should return a channel that streams text chunks.
	Interpret(divinationType string, result string, question string) (<-chan string, error)

	// ChatCompletion generates a chat response given a list of messages.
	// It should return a channel that streams text chunks.
	ChatCompletion(messages []map[string]string) (<-chan string, error)
}

// OpenAIProvider implements AIProvider using OpenAI-compatible API.
type OpenAIProvider struct {
	apiKey        string
	baseURL       string
	model         string
	fallbackModel string
	maxTokens     int
	temperature   float64
	client        *http.Client
	prompts       map[string]*template.Template
}

// NewOpenAIProvider creates a new OpenAI-compatible AI provider.
func NewOpenAIProvider(cfg *config.AIConfig) (*OpenAIProvider, error) {
	provider := &OpenAIProvider{
		apiKey:        cfg.APIKey,
		baseURL:       cfg.BaseURL,
		model:         cfg.Model,
		fallbackModel: cfg.FallbackModel,
		maxTokens:     cfg.MaxTokens,
		temperature:   cfg.Temperature,
		client: &http.Client{
			Timeout: aiHTTPTimeout,
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
	types := []string{"tarot", "horoscope", "liuyao", "liuyao_v2", "bazi"}
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
	err := tmpl.Execute(&promptBuf, buildPromptData(divinationType, result, question))
	if err != nil {
		return nil, fmt.Errorf("failed to execute prompt template: %w", err)
	}
	if divinationType == "liuyao_v2" {
		prompt := promptBuf.String()
		log.Info().
			Str("component", "ai").
			Str("divination_type", divinationType).
			Int("prompt_chars", len([]rune(prompt))).
			Int("max_tokens", p.maxTokens).
			Str("prompt", prompt).
			Msg("prepared LiuYao v2 AI prompt")
	}

	// Make streaming API request
	reqBody := map[string]interface{}{
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

	resp, err := p.doChatCompletionRequest(reqBody)
	if err != nil {
		return nil, err
	}

	// Parse streaming response
	ch := make(chan string, 100)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		var outputChars int
		var finished bool
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
				finished = true
				return
			}

			var chunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				log.Warn().
					Str("component", "ai").
					Str("divination_type", divinationType).
					Err(err).
					Str("raw_chunk", data).
					Msg("failed to parse AI stream chunk")
				continue
			}
			if len(chunk.Choices) > 0 {
				if chunk.Choices[0].Delta.Content != "" {
					outputChars += len([]rune(chunk.Choices[0].Delta.Content))
					ch <- chunk.Choices[0].Delta.Content
				}
				if chunk.Choices[0].FinishReason != "" {
					event := log.Info()
					if chunk.Choices[0].FinishReason == "length" {
						ch <- incompleteLengthNotice
						event = log.Warn()
					}
					finished = true
					event.
						Str("component", "ai").
						Str("divination_type", divinationType).
						Str("finish_reason", chunk.Choices[0].FinishReason).
						Int("output_chars", outputChars).
						Int("max_tokens", p.maxTokens).
						Msg("AI stream finished")
				}
			}
		}
		if err := scanner.Err(); err != nil {
			ch <- incompleteConnectionNotice
			log.Warn().
				Str("component", "ai").
				Str("divination_type", divinationType).
				Err(err).
				Msg("AI stream scanner stopped with error")
		} else if !finished {
			ch <- incompleteConnectionNotice
			log.Warn().
				Str("component", "ai").
				Str("divination_type", divinationType).
				Msg("AI stream ended without completion marker")
		}
	}()

	return ch, nil
}

// ChatCompletion streams a chat response given a list of messages.
func (p *OpenAIProvider) ChatCompletion(messages []map[string]string) (<-chan string, error) {
	// Make streaming API request
	reqBody := map[string]interface{}{
		"messages":    messages,
		"max_tokens":  p.maxTokens,
		"temperature": p.temperature,
		"stream":      true,
	}

	resp, err := p.doChatCompletionRequest(reqBody)
	if err != nil {
		return nil, err
	}

	// Parse streaming response
	ch := make(chan string, 100)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		var outputChars int
		var finished bool
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
				finished = true
				return
			}

			var chunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				log.Warn().
					Str("component", "ai").
					Err(err).
					Str("raw_chunk", data).
					Msg("failed to parse AI stream chunk")
				continue
			}
			if len(chunk.Choices) > 0 {
				if chunk.Choices[0].Delta.Content != "" {
					outputChars += len([]rune(chunk.Choices[0].Delta.Content))
					ch <- chunk.Choices[0].Delta.Content
				}
				if chunk.Choices[0].FinishReason != "" {
					if chunk.Choices[0].FinishReason == "length" {
						ch <- incompleteLengthNotice
					}
					finished = true
					log.Info().
						Str("component", "ai").
						Str("finish_reason", chunk.Choices[0].FinishReason).
						Int("output_chars", outputChars).
						Int("max_tokens", p.maxTokens).
						Msg("AI chat stream finished")
				}
			}
		}
		if err := scanner.Err(); err != nil {
			ch <- incompleteConnectionNotice
			log.Warn().
				Str("component", "ai").
				Err(err).
				Msg("AI chat stream scanner stopped with error")
		} else if !finished {
			ch <- incompleteConnectionNotice
			log.Warn().
				Str("component", "ai").
				Msg("AI chat stream ended without completion marker")
		}
	}()

	return ch, nil
}

func (p *OpenAIProvider) doChatCompletionRequest(reqBody map[string]interface{}) (*http.Response, error) {
	models := []string{p.model}
	if p.fallbackModel != "" && p.fallbackModel != p.model {
		models = append(models, p.fallbackModel)
	}

	var lastErr error
	for i, model := range models {
		reqBody["model"] = model
		resp, err := p.sendChatCompletionRequest(reqBody)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		lastErr = err
		if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = apperrors.New(apperrors.ErrAIServiceUnavailable,
				fmt.Sprintf("AI service returned status %d: %s", resp.StatusCode, string(body)))
		}

		if i+1 < len(models) {
			log.Warn().
				Str("component", "ai").
				Str("model", model).
				Str("fallback_model", models[i+1]).
				Err(lastErr).
				Msg("AI model request failed, retrying fallback model")
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, apperrors.New(apperrors.ErrAIServiceUnavailable, "AI service request failed")
}

func (p *OpenAIProvider) sendChatCompletionRequest(reqBody map[string]interface{}) (*http.Response, error) {
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
	return resp, nil
}

func buildPromptData(divinationType string, result string, question string) map[string]string {
	data := map[string]string{
		"Spread":       extractSpread(result),
		"Zodiac":       extractZodiac(result),
		"Period":       extractPeriod(result),
		"Question":     question,
		"Result":       result,
		"Method":       "",
		"BenGua":       "",
		"BianGua":      "无",
		"MovingLines":  "无",
		"BookEvidence": "无",
		"MethodRules":  "无",
	}
	if divinationType != "liuyao_v2" {
		return data
	}

	var reading struct {
		Method       string                       `json:"method"`
		BenGua       *model.TakashimaHexagram     `json:"ben_gua"`
		BianGua      *model.TakashimaHexagram     `json:"bian_gua"`
		MutableLines []int                        `json:"mutable_lines"`
		BookEvidence *model.TakashimaBookEvidence `json:"book_evidence"`
	}
	if err := json.Unmarshal([]byte(result), &reading); err != nil {
		return data
	}

	data["Method"] = reading.Method
	data["BenGua"] = formatPromptHexagram(reading.BenGua)
	data["BianGua"] = formatPromptHexagram(reading.BianGua)
	data["MovingLines"] = formatPromptMovingLines(reading.BenGua, reading.MutableLines)
	if reading.BookEvidence != nil {
		data["BookEvidence"] = formatEvidenceSnippets(reading.BookEvidence.Snippets)
		data["MethodRules"] = formatEvidenceSnippets(reading.BookEvidence.MethodRules)
		log.Info().
			Str("component", "ai").
			Str("divination_type", "liuyao_v2").
			Str("ben_gua", formatPromptHexagram(reading.BenGua)).
			Str("bian_gua", formatPromptHexagram(reading.BianGua)).
			Str("moving_lines", data["MovingLines"]).
			Interface("prompt_book_evidence", summarizeEvidenceForLog(reading.BookEvidence.Snippets)).
			Interface("prompt_method_rules", summarizeEvidenceForLog(reading.BookEvidence.MethodRules)).
			Msg("using Takashima book evidence in AI prompt")
	}
	return data
}

func formatMeihuaHexagram(hex model.MeiHuaHexagram) string {
	return fmt.Sprintf("%s（%s上%s下）", hex.Name, hex.Upper.Name, hex.Lower.Name)
}

func formatMovingLine(line int) string {
	names := []string{"", "初爻", "二爻", "三爻", "四爻", "五爻", "上爻"}
	if line >= 1 && line <= 6 {
		return names[line] + "动"
	}
	return "无动爻"
}

func formatPromptHexagram(hexagram *model.TakashimaHexagram) string {
	if hexagram == nil {
		return "无"
	}
	if hexagram.FullName != "" {
		return fmt.Sprintf("%s（%s上%s下）", hexagram.FullName, hexagram.UpperTrigram, hexagram.LowerTrigram)
	}
	return hexagram.Name
}

func formatPromptMovingLines(hexagram *model.TakashimaHexagram, mutableLines []int) string {
	if hexagram == nil || len(mutableLines) == 0 {
		return "无动爻，以本卦卦辞与卦象为主。"
	}
	parts := make([]string, 0, len(mutableLines))
	for _, idx := range mutableLines {
		position := idx + 1
		for _, line := range hexagram.Lines {
			if line.Position == position {
				parts = append(parts, fmt.Sprintf("%s：%s", line.Name, line.Original))
				break
			}
		}
	}
	if len(parts) == 0 {
		return "动爻未匹配到书中爻辞。"
	}
	return strings.Join(parts, "\n")
}

func formatEvidenceSnippets(snippets []model.TakashimaEvidenceSnippet) string {
	if len(snippets) == 0 {
		return "无"
	}
	var b strings.Builder
	for i, snippet := range snippets {
		if i > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString("【")
		b.WriteString(snippet.Title)
		b.WriteString("】\n")
		b.WriteString(snippet.Text)
	}
	return b.String()
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
	profiles *config.ProfilesConfig
}

// NewAIService creates a new AIService.
func NewAIService(provider AIProvider, reader AIResultReader, profiles *config.ProfilesConfig) *AIService {
	return &AIService{
		provider: provider,
		reader:   reader,
		profiles: profiles,
	}
}

// Interpret starts streaming AI interpretation for a divination record.
func (s *AIService) Interpret(userID uint, recordID uint, divinationType string, question string, force bool) (<-chan string, error) {
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
	if record.AIReading != "" && !force {
		// Return existing reading as a single chunk
		ch := make(chan string, 1)
		ch <- record.AIReading
		close(ch)
		return ch, nil
	}

	// Start streaming interpretation
	ch, err := s.streamInterpretation(divinationType, record.Result, question)
	if err != nil {
		return nil, err
	}

	// Collect the full reading in background for saving
	fullReadingCh := make(chan string, 100)
	go func() {
		defer close(fullReadingCh)

		var fullReading strings.Builder
		for chunk := range ch {
			fullReading.WriteString(chunk)
			fullReadingCh <- chunk
		}

		// Save only complete readings, so history does not cache truncated output.
		reading := fullReading.String()
		if reading != "" && !isIncompleteAIReading(reading) {
			_ = s.reader.UpdateAIReading(recordID, reading)
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

	return s.streamInterpretation(divinationType, result, question)
}

// streamInterpretation routes to Compose-based or legacy prompt based on divination type.
func (s *AIService) streamInterpretation(divinationType string, result string, question string) (<-chan string, error) {
	// For meihua, use Compose with profile
	if divinationType == "meihua" {
		log.Info().
			Str("divination_type", divinationType).
			Bool("has_profiles", s.profiles != nil).
			Str("question", truncateForLog(question, 60)).
			Msg("streamInterpretation: routing meihua via Compose pipeline")

		if s.profiles == nil {
			log.Warn().Msg("streamInterpretation: profiles config is nil, falling back to legacy template")
			return s.provider.Interpret(divinationType, result, question)
		}

		// Step 1: Resolve default profile
		profile, err := s.profiles.DefaultProfile("meihua")
		if err != nil {
			log.Warn().
				Err(err).
				Str("divination_type", "meihua").
				Msg("streamInterpretation: DefaultProfile failed, falling back to legacy template")
			return s.provider.Interpret(divinationType, result, question)
		}
		log.Info().
			Str("profile_id", s.profiles.DefaultBindings["meihua"]).
			Str("profile_name", profile.Name).
			Str("profile_version", profile.Version).
			Msg("streamInterpretation: profile resolved ✓")

		// Step 2: Parse facts from result JSON
		facts, err := ParseFacts("meihua", result, question)
		if err != nil {
			log.Warn().
				Err(err).
				Msg("streamInterpretation: ParseFacts failed, falling back to legacy template")
			return s.provider.Interpret(divinationType, result, question)
		}
		log.Info().
			Int("fact_count", len(facts)).
			Strs("fact_keys", factKeysList(facts)).
			Msg("streamInterpretation: facts parsed ✓")

		// Step 3: Compose system + user messages
		messages, err := Compose(profile, InterpretationInput{
			DivinationType: "meihua",
			Question:       question,
			ResultJSON:     result,
			ResultFacts:    facts,
			Mode:           "direct",
		})
		if err != nil {
			log.Warn().
				Err(err).
				Msg("streamInterpretation: Compose failed, falling back to legacy template")
			return s.provider.Interpret(divinationType, result, question)
		}
		log.Info().
			Int("system_chars", len([]rune(messages.System))).
			Int("user_chars", len([]rune(messages.User))).
			Str("mode", "direct").
			Msg("streamInterpretation: Compose succeeded ✓")

		return s.provider.ChatCompletion([]map[string]string{
			{"role": "system", "content": messages.System},
			{"role": "user", "content": messages.User},
		})
	}

	// Fallback: legacy template-based prompt
	log.Debug().
		Str("divination_type", divinationType).
		Msg("streamInterpretation: using legacy template prompt (no Compose route defined)")
	return s.provider.Interpret(divinationType, result, question)
}

// truncateForLog truncates a string for log display.
func truncateForLog(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

// factKeysList returns the keys of a facts map as a sorted slice for logging.
func factKeysList(facts map[string]string) []string {
	keys := make([]string, 0, len(facts))
	for k := range facts {
		keys = append(keys, k)
	}
	return keys
}

func isIncompleteAIReading(reading string) bool {
	return strings.Contains(reading, "AI 输出达到长度上限") ||
		strings.Contains(reading, "AI 解读连接中断") ||
		strings.Contains(reading, "内容可能未完整生成")
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

// ChatCompletion returns a mock chat response.
func (m *MockAIProvider) ChatCompletion(messages []map[string]string) (<-chan string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	ch := make(chan string, 100)
	go func() {
		defer close(ch)
		if m.Response != "" {
			ch <- m.Response
		}
	}()
	return ch, nil
}
