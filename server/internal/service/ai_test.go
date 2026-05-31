package service

import (
	"testing"

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

	svc := NewAIService(provider, reader)
	assert.NotNil(t, svc)
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

	svc := NewAIService(provider, reader)

	ch, err := svc.Interpret(1, 1, "tarot", "")
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

func TestAIService_Interpret_RecordNotFound(t *testing.T) {
	provider := &MockAIProvider{Response: "test"}
	reader := NewMockResultReader()

	svc := NewAIService(provider, reader)

	_, err := svc.Interpret(1, 999, "tarot", "")
	assert.Error(t, err)
}

func TestAIService_InterpretDirect_ProviderUnavailable(t *testing.T) {
	svc := NewAIService(nil, NewMockResultReader())

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

	svc := NewAIService(provider, reader)

	ch, err := svc.Interpret(1, 1, "tarot", "")
	require.NoError(t, err)

	// Should get the existing reading
	var result string
	for chunk := range ch {
		result += chunk
	}

	assert.Equal(t, "existing reading", result)
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

	svc := NewAIService(provider, reader)

	// Test with custom question
	ch, err := svc.Interpret(1, 1, "tarot", "自定义问题")
	require.NoError(t, err)

	var result string
	for chunk := range ch {
		result += chunk
	}

	assert.NotEmpty(t, result)
}

