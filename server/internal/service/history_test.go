package service

import (
	"encoding/json"
	"testing"
	"time"

	"gorm.io/gorm"
	"zhanbu/internal/model"
)

// mockDivinationRepo is a simple in-memory mock for testing.
type mockDivinationRepo struct {
	records  []model.DivinationRecord
	nextID   uint
	findErr  error
	createErr error
	deleteErr error
}

func newMockDivinationRepo() *mockDivinationRepo {
	return &mockDivinationRepo{nextID: 1, findErr: gorm.ErrRecordNotFound}
}

func (m *mockDivinationRepo) Create(record *model.DivinationRecord) error {
	if m.createErr != nil {
		return m.createErr
	}
	record.ID = m.nextID
	m.nextID++
	record.CreatedAt = time.Now()
	m.records = append(m.records, *record)
	return nil
}

func (m *mockDivinationRepo) FindByID(id uint) (*model.DivinationRecord, error) {
	for _, r := range m.records {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, m.findErr
}

func (m *mockDivinationRepo) FindByUserIDAndID(userID uint, id uint) (*model.DivinationRecord, error) {
	for _, r := range m.records {
		if r.ID == id && r.UserID == userID {
			return &r, nil
		}
	}
	return nil, m.findErr
}

func (m *mockDivinationRepo) ListByUserID(userID uint, divinationType string, page, pageSize int) ([]model.DivinationRecord, int64, error) {
	var filtered []model.DivinationRecord
	for _, r := range m.records {
		if r.UserID == userID {
			if divinationType == "" || r.Type == divinationType {
				filtered = append(filtered, r)
			}
		}
	}
	total := int64(len(filtered))
	start := (page - 1) * pageSize
	if start > len(filtered) {
		return []model.DivinationRecord{}, total, nil
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (m *mockDivinationRepo) DeleteByUserIDAndID(userID uint, id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	for i, r := range m.records {
		if r.ID == id && r.UserID == userID {
			m.records = append(m.records[:i], m.records[i+1:]...)
			return nil
		}
	}
	return m.findErr // record not found
}

// seedRecords creates test data in the mock repo.
func seedRecords(repo *mockDivinationRepo) {
	resultJSON, _ := json.Marshal(map[string]interface{}{
		"spread": "single",
		"cards": []map[string]interface{}{
			{"card": map[string]interface{}{"name": "愚者"}, "orientation": "upright"},
		},
	})

	repo.Create(&model.DivinationRecord{UserID: 1, Type: "tarot", Question: "今日运势", Result: string(resultJSON)})
	repo.Create(&model.DivinationRecord{UserID: 1, Type: "tarot", Question: "工作如何", Result: string(resultJSON), AIReading: "工作运势不错，适合开展新项目"})
	repo.Create(&model.DivinationRecord{UserID: 1, Type: "bazi", Question: "婚姻", Result: `{"year":"甲子"}`})
	repo.Create(&model.DivinationRecord{UserID: 2, Type: "tarot", Question: "别人的问题", Result: string(resultJSON)})
}

func TestHistoryList_All(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	result, appErr := svc.List(1, "", 1, 10)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if result.Total != 3 {
		t.Errorf("expected total 3, got %d", result.Total)
	}
	if len(result.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(result.Items))
	}
	if result.Page != 1 {
		t.Errorf("expected page 1, got %d", result.Page)
	}
	if result.PageSize != 10 {
		t.Errorf("expected page_size 10, got %d", result.PageSize)
	}
}

func TestHistoryList_FilterByType(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	result, appErr := svc.List(1, "tarot", 1, 10)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if result.Total != 2 {
		t.Errorf("expected total 2 for tarot, got %d", result.Total)
	}
	for _, item := range result.Items {
		if item.Type != "tarot" {
			t.Errorf("expected type 'tarot', got '%s'", item.Type)
		}
	}
}

func TestHistoryList_FilterByBazi(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	result, appErr := svc.List(1, "bazi", 1, 10)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1 for bazi, got %d", result.Total)
	}
}

func TestHistoryList_Pagination(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	// Page 1 with size 2
	result, appErr := svc.List(1, "", 1, 2)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if result.Total != 3 {
		t.Errorf("expected total 3, got %d", result.Total)
	}
	if len(result.Items) != 2 {
		t.Errorf("expected 2 items on page 1, got %d", len(result.Items))
	}

	// Page 2 with size 2
	result, appErr = svc.List(1, "", 2, 2)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if len(result.Items) != 1 {
		t.Errorf("expected 1 item on page 2, got %d", len(result.Items))
	}
}

func TestHistoryList_OtherUserIsolation(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	// User 2 should only see their own record
	result, appErr := svc.List(2, "", 1, 10)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1 for user 2, got %d", result.Total)
	}
	if result.Items[0].Question != "别人的问题" {
		t.Errorf("expected '别人的问题', got '%s'", result.Items[0].Question)
	}
}

func TestHistoryList_TypeCNMapping(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	result, appErr := svc.List(1, "tarot", 1, 10)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	for _, item := range result.Items {
		if item.TypeCN != "塔罗牌" {
			t.Errorf("expected type_cn '塔罗牌', got '%s'", item.TypeCN)
		}
	}
}

func TestHistoryList_AIReadingSummary(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	result, appErr := svc.List(1, "tarot", 1, 10)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}

	// The second tarot record has AIReading
	found := false
	for _, item := range result.Items {
		if item.Question == "工作如何" {
			found = true
			if item.Summary != "工作运势不错，适合开展新项目" {
				t.Errorf("expected AI reading as summary, got '%s'", item.Summary)
			}
		}
	}
	if !found {
		t.Error("expected to find '工作如何' item")
	}
}

func TestHistoryList_DefaultPagination(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	// page < 1 should default to 1, pageSize < 1 should default to 10
	result, appErr := svc.List(1, "", 0, 0)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if result.Page != 1 {
		t.Errorf("expected page 1, got %d", result.Page)
	}
	if result.PageSize != 10 {
		t.Errorf("expected page_size 10, got %d", result.PageSize)
	}
}

func TestGetDetail_Success(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	record, appErr := svc.GetDetail(1, 1)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if record.UserID != 1 {
		t.Errorf("expected user_id 1, got %d", record.UserID)
	}
	if record.Question != "今日运势" {
		t.Errorf("expected question '今日运势', got '%s'", record.Question)
	}
}

func TestGetDetail_NotFound(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	_, appErr := svc.GetDetail(1, 999)
	if appErr == nil {
		t.Fatal("expected error, got nil")
	}
	if appErr.Code != 1004 {
		t.Errorf("expected error code 1004 (not found), got %d", appErr.Code)
	}
}

func TestGetDetail_WrongUser(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	// User 2 tries to access user 1's record
	_, appErr := svc.GetDetail(2, 1)
	if appErr == nil {
		t.Fatal("expected error, got nil")
	}
	if appErr.Code != 1004 {
		t.Errorf("expected error code 1004 (not found), got %d", appErr.Code)
	}
}

func TestDelete_Success(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	appErr := svc.Delete(1, 1)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}

	// Verify deleted
	_, appErr = svc.GetDetail(1, 1)
	if appErr == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestDelete_NotFound(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	appErr := svc.Delete(1, 999)
	if appErr == nil {
		t.Fatal("expected error, got nil")
	}
	if appErr.Code != 1004 {
		t.Errorf("expected error code 1004 (not found), got %d", appErr.Code)
	}
}

func TestDelete_WrongUser(t *testing.T) {
	repo := newMockDivinationRepo()
	seedRecords(repo)
	svc := NewHistoryServiceWithReader(repo)

	// User 2 tries to delete user 1's record
	appErr := svc.Delete(2, 1)
	if appErr == nil {
		t.Fatal("expected error, got nil")
	}
	if appErr.Code != 1004 {
		t.Errorf("expected error code 1004 (not found), got %d", appErr.Code)
	}

	// Verify record still exists for user 1
	_, appErr = svc.GetDetail(1, 1)
	if appErr != nil {
		t.Errorf("record should still exist for user 1: %v", appErr)
	}
}

func TestHistoryList_EmptyResult(t *testing.T) {
	repo := newMockDivinationRepo()
	// No records seeded
	svc := NewHistoryServiceWithReader(repo)

	result, appErr := svc.List(1, "", 1, 10)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if result.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Total)
	}
	if len(result.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Items))
	}
}

func TestBuildSummary_TarotCards(t *testing.T) {
	resultJSON := `{"spread":"three","cards":[
		{"card":{"name":"愚者"},"orientation":"upright"},
		{"card":{"name":"魔术师"},"orientation":"reversed"},
		{"card":{"name":"女祭司"},"orientation":"upright"}
	]}`
	summary := buildSummary(resultJSON)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	// Should contain at least first card name
	if len(summary) < 2 {
		t.Errorf("summary too short: '%s'", summary)
	}
}

func TestBuildSummary_InvalidJSON(t *testing.T) {
	summary := buildSummary("not json at all and this is a very long string that exceeds one hundred characters for sure yes it does and more and more and more padding text here okay done")
	if len(summary) > 103 { // 100 + "..."
		t.Errorf("summary should be truncated, got length %d", len(summary))
	}
}

func TestBuildSummary_ShortText(t *testing.T) {
	summary := buildSummary("short")
	if summary != "short" {
		t.Errorf("expected 'short', got '%s'", summary)
	}
}
