package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zhanbu/internal/model"
)

type mockChatRepo struct {
	nextSessionID uint
	nextMessageID uint
	sessions      map[uint]*model.ChatSession
	messages      map[uint][]model.ChatMessage
}

func newMockChatRepo() *mockChatRepo {
	return &mockChatRepo{
		nextSessionID: 1,
		nextMessageID: 1,
		sessions:      make(map[uint]*model.ChatSession),
		messages:      make(map[uint][]model.ChatMessage),
	}
}

func (m *mockChatRepo) CreateSession(session *model.ChatSession) error {
	session.ID = m.nextSessionID
	m.nextSessionID++
	session.CreatedAt = time.Now()
	session.UpdatedAt = session.CreatedAt
	copied := *session
	m.sessions[session.ID] = &copied
	return nil
}

func (m *mockChatRepo) FindSessionByID(id uint) (*model.ChatSession, error) {
	session, ok := m.sessions[id]
	if !ok {
		return nil, errors.New("session not found")
	}
	copied := *session
	copied.Messages = append([]model.ChatMessage(nil), m.messages[id]...)
	return &copied, nil
}

func (m *mockChatRepo) FindSessionByUserAndID(userID uint, sessionID uint) (*model.ChatSession, error) {
	session, err := m.FindSessionByID(sessionID)
	if err != nil || session.UserID != userID {
		return nil, errors.New("session not found")
	}
	return session, nil
}

func (m *mockChatRepo) FindSessionByUserAndRecord(userID uint, recordID uint) (*model.ChatSession, error) {
	for _, session := range m.sessions {
		if session.UserID == userID && session.RecordID == recordID {
			return session, nil
		}
	}
	return nil, errors.New("session not found")
}

func (m *mockChatRepo) ListSessionsByUser(userID uint, page, size int) ([]model.ChatSession, int64, error) {
	return nil, 0, nil
}

func (m *mockChatRepo) UpdateSessionTitle(id uint, title string) error {
	if session, ok := m.sessions[id]; ok {
		session.Title = title
	}
	return nil
}

func (m *mockChatRepo) UpdateSessionTimestamp(id uint) error {
	if session, ok := m.sessions[id]; ok {
		session.UpdatedAt = time.Now()
	}
	return nil
}

func (m *mockChatRepo) DeleteSession(id uint) error {
	delete(m.sessions, id)
	delete(m.messages, id)
	return nil
}

func (m *mockChatRepo) CreateMessage(message *model.ChatMessage) error {
	message.ID = m.nextMessageID
	m.nextMessageID++
	message.CreatedAt = time.Now()
	m.messages[message.SessionID] = append(m.messages[message.SessionID], *message)
	return nil
}

func (m *mockChatRepo) GetRecentMessages(sessionID uint, limit int) ([]model.ChatMessage, error) {
	messages := append([]model.ChatMessage(nil), m.messages[sessionID]...)
	if len(messages) > limit {
		messages = messages[len(messages)-limit:]
	}
	return messages, nil
}

func (m *mockChatRepo) CountMessagesBySession(sessionID uint) (int64, error) {
	return int64(len(m.messages[sessionID])), nil
}

type fakeChatStarter struct {
	record *model.DivinationRecord
}

func (s *fakeChatStarter) Start(userID uint, divinationType string, question string) (*model.DivinationRecord, error) {
	s.record = &model.DivinationRecord{
		ID:       42,
		UserID:   userID,
		Type:     divinationType,
		Question: question,
		Result:   `{"mode":"test"}`,
	}
	return s.record, nil
}

func TestChatServiceCreateModeSessionCreatesRecordAndUserMessage(t *testing.T) {
	chatRepo := newMockChatRepo()
	recordRepo := newMockDivinationRepo()
	recordRepo.records = append(recordRepo.records, model.DivinationRecord{
		ID:       42,
		UserID:   7,
		Type:     "tarot",
		Question: "我今天适合做决定吗",
		Result:   `{"mode":"test"}`,
	})
	starter := &fakeChatStarter{}
	svc := NewChatService(chatRepo, recordRepo, &MockAIProvider{Response: "这是一段聊天模式AI解读"})
	svc.SetDivinationStarter(starter)

	session, err := svc.CreateModeSession(7, ChatStartOptions{
		Type:     "tarot",
		Question: "我今天适合做决定吗",
	})

	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, uint(42), session.RecordID)
	require.Len(t, session.Messages, 1)
	assert.Equal(t, "user", session.Messages[0].Role)
	assert.Equal(t, "我今天适合做决定吗", session.Messages[0].Content)
	assert.Empty(t, recordRepo.records[0].AIReading)
}

type controllableChatProvider struct {
	ch chan string
}

func (p *controllableChatProvider) Interpret(divinationType string, result string, question string) (<-chan string, error) {
	return nil, nil
}

func (p *controllableChatProvider) ChatCompletion(messages []map[string]string) (<-chan string, error) {
	return p.ch, nil
}

type failingChatProvider struct{}

func (p *failingChatProvider) Interpret(divinationType string, result string, question string) (<-chan string, error) {
	return nil, nil
}

func (p *failingChatProvider) ChatCompletion(messages []map[string]string) (<-chan string, error) {
	return nil, errors.New("chat completion should not be called")
}

func TestChatServiceCreateModeSessionDoesNotWaitForAIReading(t *testing.T) {
	chatRepo := newMockChatRepo()
	recordRepo := newMockDivinationRepo()
	recordRepo.records = append(recordRepo.records, model.DivinationRecord{
		ID:       42,
		UserID:   7,
		Type:     "liuyao_v2",
		Question: "今天晚上是否可以出门",
		Result:   `{"mode":"test"}`,
	})
	starter := &fakeChatStarter{}
	svc := NewChatService(chatRepo, recordRepo, &failingChatProvider{})
	svc.SetDivinationStarter(starter)

	session, err := svc.CreateModeSession(7, ChatStartOptions{
		Type:     "liuyao_v2",
		Question: "今天晚上是否可以出门",
	})

	require.NoError(t, err)
	require.NotNil(t, session)
	require.Len(t, session.Messages, 1)
	assert.Equal(t, "user", session.Messages[0].Role)
	assert.Empty(t, recordRepo.records[0].AIReading)
}

func TestChatServiceStreamInitialReadingSavesAIReading(t *testing.T) {
	chatRepo := newMockChatRepo()
	recordRepo := newMockDivinationRepo()
	recordRepo.records = append(recordRepo.records, model.DivinationRecord{
		ID:       42,
		UserID:   7,
		Type:     "tarot",
		Question: "我今天适合做决定吗",
		Result:   `{"mode":"test"}`,
	})
	require.NoError(t, chatRepo.CreateSession(&model.ChatSession{
		UserID:   7,
		RecordID: 42,
		Title:    "我今天适合做决定吗",
	}))
	require.NoError(t, chatRepo.CreateMessage(&model.ChatMessage{
		SessionID: 1,
		Role:      "user",
		Content:   "我今天适合做决定吗",
	}))
	svc := NewChatService(chatRepo, recordRepo, &MockAIProvider{Response: "这是一段聊天模式AI解读"})

	stream, err := svc.StreamInitialReading(7, 1)
	require.NoError(t, err)

	var got string
	for chunk := range stream {
		got += chunk
	}

	assert.Equal(t, "这是一段聊天模式AI解读", got)
	require.Len(t, chatRepo.messages[1], 2)
	assert.Equal(t, "assistant", chatRepo.messages[1][1].Role)
	assert.Equal(t, "这是一段聊天模式AI解读", chatRepo.messages[1][1].Content)
	assert.Equal(t, "这是一段聊天模式AI解读", recordRepo.records[0].AIReading)
}

func TestFormatMeihuaMethodDescPrefersChineseLunarDisplay(t *testing.T) {
	resultJSON := `{"source_values":{"year_branch":"午","lunar_month":5,"lunar_day":9,"lunar_display":"丙午年五月初九子时","hour_branch":"子"}}`

	got := formatMeihuaMethodDesc("time", resultJSON)

	assert.Equal(t, "时间起卦（丙午年五月初九子时）", got)
}

func TestChatServiceSendMessageBroadcastsFullStreamAndSavesFullResponse(t *testing.T) {
	chatRepo := newMockChatRepo()
	recordRepo := newMockDivinationRepo()
	recordRepo.records = append(recordRepo.records, model.DivinationRecord{
		ID:       12,
		UserID:   7,
		Type:     "liuyao_v2",
		Question: "今晚适合出门吗",
		Result:   `{"mode":"test"}`,
	})
	require.NoError(t, chatRepo.CreateSession(&model.ChatSession{
		UserID:   7,
		RecordID: 12,
		Title:    "今晚适合出门吗",
	}))
	provider := &controllableChatProvider{ch: make(chan string)}
	svc := NewChatService(chatRepo, recordRepo, provider)

	stream, err := svc.SendMessage(7, 1, "继续说")
	require.NoError(t, err)

	done := make(chan string, 1)
	go func() {
		var got string
		for chunk := range stream {
			got += chunk
		}
		done <- got
	}()

	provider.ch <- "第一段"
	provider.ch <- "第二段"
	provider.ch <- "第三段"
	close(provider.ch)

	assert.Equal(t, "第一段第二段第三段", <-done)
	require.Eventually(t, func() bool {
		messages := chatRepo.messages[1]
		return len(messages) == 2 && messages[1].Role == "assistant"
	}, time.Second, 10*time.Millisecond)
	assert.Equal(t, "第一段第二段第三段", chatRepo.messages[1][1].Content)
}
