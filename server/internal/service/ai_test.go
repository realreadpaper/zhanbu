package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zhanbu/internal/model"
	apperrors "zhanbu/pkg/errors"
)

// MockResultReader is a mock implementation of AIResultReader for testing.
type MockResultReader struct {
	records map[uint]*model.DivinationRecord
	reading map[uint]string
}

// NewMockResultReader creates a new MockResultReader.
func NewMockResultReader() *MockResultReader {
	return &MockResultReader{
		records: make(map[uint]*model.DivinationRecord),
		reading: make(map[uint]string),
	}
}

// FindByUserIDAndID returns a mock record.
func (m *MockResultReader) FindByUserIDAndID(userID uint, id uint) (*model.DivinationRecord, error) {
	record, ok := m.records[id]
	if !ok {
		return nil, assert.AnError
	}
	if record.UserID != userID {
		return nil, assert.AnError
	}
	return record, nil
}

// UpdateAIReading saves the reading.
func (m *MockResultReader) UpdateAIReading(id uint, reading string) error {
	m.reading[id] = reading
	return nil
}

func TestNewAIService(t *testing.T) {
	provider := &MockAIProvider{Response: "test"}
	reader := NewMockResultReader()

	svc := NewAIService(provider, reader, nil)
	assert.NotNil(t, svc)
}

func TestOpenAIProviderLoadPromptsDoesNotRegisterMeihuaLegacyPrompt(t *testing.T) {
	provider := &OpenAIProvider{
		prompts: make(map[string]*template.Template),
	}

	err := provider.loadPrompts()

	require.NoError(t, err)
	assert.NotContains(t, provider.prompts, "meihua")
	assert.Contains(t, provider.prompts, "liuyao_v2")
}

func TestOpenAIProviderChatCompletionFallsBackToFlashModel(t *testing.T) {
	var requestedModels []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Model string `json:"model"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		requestedModels = append(requestedModels, payload.Model)

		if payload.Model == "deepseek-v4-pro" {
			http.Error(w, "pro unavailable", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"ok\"}}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer server.Close()

	provider := &OpenAIProvider{
		apiKey:        "sk-test",
		baseURL:       server.URL,
		model:         "deepseek-v4-pro",
		fallbackModel: "deepseek-v4-flash",
		maxTokens:     1000,
		temperature:   0.7,
		client:        server.Client(),
		prompts:       make(map[string]*template.Template),
	}

	ch, err := provider.ChatCompletion([]map[string]string{{"role": "user", "content": "ping"}})

	require.NoError(t, err)
	var got string
	for chunk := range ch {
		got += chunk
	}
	assert.Equal(t, "ok", got)
	assert.Equal(t, []string{"deepseek-v4-pro", "deepseek-v4-flash"}, requestedModels)
}

func TestAIService_Interpret_Success(t *testing.T) {
	provider := &MockAIProvider{Response: "这是一个测试解读"}
	reader := NewMockResultReader()

	// Add a test record
	reader.records[1] = &model.DivinationRecord{
		ID:       1,
		UserID:   1,
		Type:     "tarot",
		Question: "我的事业如何？",
		Result:   `{"spread":"single","cards":[],"timestamp":"2026-01-01"}`,
	}

	svc := NewAIService(provider, reader, nil)

	ch, err := svc.Interpret(1, 1, "tarot", "", false)
	require.NoError(t, err)
	require.NotNil(t, ch)

	// Collect all chunks
	var fullReading string
	for chunk := range ch {
		fullReading += chunk
	}

	assert.NotEmpty(t, fullReading)
	assert.Equal(t, "这是一个测试解读", reader.reading[1])
}

func TestAIService_Interpret_DoesNotSaveIncompleteReading(t *testing.T) {
	provider := &MockAIProvider{Response: "前半段解读\n\n【系统提示：AI 输出达到长度上限，请稍后重新解读。】"}
	reader := NewMockResultReader()

	reader.records[1] = &model.DivinationRecord{
		ID:       1,
		UserID:   1,
		Type:     "liuyao_v2",
		Question: "关系能长久吗",
		Result:   `{"method":"coins","ben_gua":{"name":"乾"}}`,
	}

	svc := NewAIService(provider, reader, nil)

	ch, err := svc.Interpret(1, 1, "liuyao_v2", "", false)
	require.NoError(t, err)

	var fullReading string
	for chunk := range ch {
		fullReading += chunk
	}

	assert.Contains(t, fullReading, "AI 输出达到长度上限")
	assert.Empty(t, reader.reading[1])
}

func TestAIService_Interpret_RecordNotFound(t *testing.T) {
	provider := &MockAIProvider{Response: "test"}
	reader := NewMockResultReader()

	svc := NewAIService(provider, reader, nil)

	_, err := svc.Interpret(1, 999, "tarot", "", false)
	assert.Error(t, err)
}

func TestAIService_InterpretDirect_ProviderUnavailable(t *testing.T) {
	svc := NewAIService(nil, NewMockResultReader(), nil)

	ch, err := svc.InterpretDirect("bazi", `{"birth":{"solar":"2026-05-31 16:08"}}`, "八字排盘")

	require.Error(t, err)
	assert.Nil(t, ch)

	var appErr *apperrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, apperrors.ErrAIServiceUnavailable, appErr.Code)
}

func TestAIService_Interpret_ExistingReading(t *testing.T) {
	provider := &MockAIProvider{Response: "new reading"}
	reader := NewMockResultReader()

	// Add a record with existing reading
	reader.records[1] = &model.DivinationRecord{
		ID:        1,
		UserID:    1,
		Type:      "tarot",
		Question:  "test",
		Result:    "{}",
		AIReading: "existing reading",
	}

	svc := NewAIService(provider, reader, nil)

	ch, err := svc.Interpret(1, 1, "tarot", "", false)
	require.NoError(t, err)

	// Should get the existing reading
	var result string
	for chunk := range ch {
		result += chunk
	}

	assert.Equal(t, "existing reading", result)
}

func TestAIService_Interpret_ForceRegeneratesExistingReading(t *testing.T) {
	provider := &MockAIProvider{Response: "new reading"}
	reader := NewMockResultReader()

	reader.records[1] = &model.DivinationRecord{
		ID:        1,
		UserID:    1,
		Type:      "liuyao_v2",
		Question:  "test",
		Result:    "{}",
		AIReading: "old partial reading",
	}

	svc := NewAIService(provider, reader, nil)

	ch, err := svc.Interpret(1, 1, "liuyao_v2", "", true)
	require.NoError(t, err)

	var result string
	for chunk := range ch {
		result += chunk
	}

	assert.Equal(t, "new reading", result)
	assert.Equal(t, "new reading", reader.reading[1])
}

func TestMockAIProvider_Interpret(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		err        error
		wantChunks int
		wantErr    bool
	}{
		{
			name:       "successful streaming",
			response:   "测试",
			err:        nil,
			wantChunks: 1,
			wantErr:    false,
		},
		{
			name:       "error",
			response:   "",
			err:        assert.AnError,
			wantChunks: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &MockAIProvider{
				Response: tt.response,
				Err:      tt.err,
			}

			ch, err := provider.Interpret("tarot", "{}", "test")
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, ch)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, ch)

			var chunks []string
			for chunk := range ch {
				chunks = append(chunks, chunk)
			}

			assert.Equal(t, tt.wantChunks, len(chunks))
		})
	}
}

func TestExtractFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		function string
		expected string
	}{
		{
			name:     "extract spread",
			input:    `{"spread":"celtic","cards":[]}`,
			function: "spread",
			expected: "celtic",
		},
		{
			name:     "extract zodiac",
			input:    `{"zodiac":"aries","period":"daily"}`,
			function: "zodiac",
			expected: "aries",
		},
		{
			name:     "extract period",
			input:    `{"zodiac":"aries","period":"daily"}`,
			function: "period",
			expected: "daily",
		},
		{
			name:     "invalid json",
			input:    `invalid`,
			function: "spread",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			switch tt.function {
			case "spread":
				result = extractSpread(tt.input)
			case "zodiac":
				result = extractZodiac(tt.input)
			case "period":
				result = extractPeriod(tt.input)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAIService_Interpret_WithQuestion(t *testing.T) {
	provider := &MockAIProvider{Response: "解读结果"}
	reader := NewMockResultReader()

	reader.records[1] = &model.DivinationRecord{
		ID:       1,
		UserID:   1,
		Type:     "tarot",
		Question: "原问题",
		Result:   `{"spread":"single","cards":[]}`,
	}

	svc := NewAIService(provider, reader, nil)

	// Test with custom question
	ch, err := svc.Interpret(1, 1, "tarot", "自定义问题", false)
	require.NoError(t, err)

	var result string
	for chunk := range ch {
		result += chunk
	}

	assert.NotEmpty(t, result)
}
