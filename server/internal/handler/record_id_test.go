package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"zhanbu/internal/model"
	"zhanbu/internal/service"
)

type savedRecordRepo struct {
	nextID  uint
	records []*model.DivinationRecord
}

func (r *savedRecordRepo) Create(record *model.DivinationRecord) error {
	if r.nextID == 0 {
		r.nextID = 1
	}
	record.ID = r.nextID
	record.CreatedAt = time.Now()
	r.nextID++
	r.records = append(r.records, record)
	return nil
}

func performAuthenticatedJSON(handler gin.HandlerFunc, body string) (int, map[string]interface{}) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("user_id", uint(42))

	handler(ctx)

	var payload map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &payload)
	return rec.Code, payload
}

func dataMap(t *testing.T, payload map[string]interface{}) map[string]interface{} {
	t.Helper()
	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected object data, got %#v", payload["data"])
	}
	return data
}

func assertRecordID(t *testing.T, data map[string]interface{}, expected uint) {
	t.Helper()
	got, ok := data["record_id"].(float64)
	if !ok {
		t.Fatalf("expected record_id in response data, got %#v", data)
	}
	if uint(got) != expected {
		t.Fatalf("expected record_id %d, got %v", expected, got)
	}
}

func TestLiuYaoThrowReturnsSavedRecordID(t *testing.T) {
	repo := &savedRecordRepo{}
	handler := NewLiuYaoHandler(service.NewLiuYaoService(), repo)

	status, payload := performAuthenticatedJSON(handler.Throw, `{"question":"事业如何"}`)

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %#v", status, payload)
	}
	if len(repo.records) != 1 {
		t.Fatalf("expected one saved record, got %d", len(repo.records))
	}
	assertRecordID(t, dataMap(t, payload), repo.records[0].ID)
}

func TestLiuYaoV2ThrowReturnsSavedRecordID(t *testing.T) {
	repo := &savedRecordRepo{}
	svc, err := service.NewLiuYaoV2Service("yarrow")
	if err != nil {
		t.Fatalf("new liuyao v2 service: %v", err)
	}
	handler := NewLiuYaoV2Handler(svc, repo)

	status, payload := performAuthenticatedJSON(handler.Throw, `{"question":"感情如何","method":"coin"}`)

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %#v", status, payload)
	}
	if len(repo.records) != 1 {
		t.Fatalf("expected one saved record, got %d", len(repo.records))
	}
	assertRecordID(t, dataMap(t, payload), repo.records[0].ID)
}

func TestBaZiCalculateReturnsSavedRecordID(t *testing.T) {
	repo := &savedRecordRepo{}
	handler := NewBaZiHandler(service.NewBaZiService(), repo)

	status, payload := performAuthenticatedJSON(handler.Calculate, `{"birth_date":"1990-01-01","birth_time":"08:30","gender":"male"}`)

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %#v", status, payload)
	}
	if len(repo.records) != 1 {
		t.Fatalf("expected one saved record, got %d", len(repo.records))
	}
	assertRecordID(t, dataMap(t, payload), repo.records[0].ID)
}

type horoscopeTemplateStub struct{}

func (horoscopeTemplateStub) GetTemplate(_, _, dimension, _ string) string {
	return dimension + " text"
}

func (horoscopeTemplateStub) GetSummary(_, _, _ string) string {
	return "summary text"
}

func TestHoroscopeReturnsSavedRecordID(t *testing.T) {
	repo := &savedRecordRepo{}
	handler := NewHoroscopeHandler(service.NewHoroscopeService(horoscopeTemplateStub{}), repo)

	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/horoscope/aries?period=daily&date=2026-06-07", nil)
	ctx.Params = gin.Params{{Key: "zodiac", Value: "aries"}}
	ctx.Set("user_id", uint(42))

	handler.GetHoroscope(ctx)

	var payload map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &payload)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %#v", rec.Code, payload)
	}
	if len(repo.records) != 1 {
		t.Fatalf("expected one saved record, got %d", len(repo.records))
	}
	assertRecordID(t, dataMap(t, payload), repo.records[0].ID)
}
