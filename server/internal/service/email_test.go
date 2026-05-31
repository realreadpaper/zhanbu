package service

import (
	"testing"
	"time"

	"gorm.io/gorm"
	"zhanbu/config"
	"zhanbu/internal/model"
)

// mockVerificationRepo is a mock for testing.
type mockVerificationRepo struct {
	verifications []model.EmailVerification
	sendLogs      []model.SendLog
	nextID        uint
}

func (m *mockVerificationRepo) Create(v *model.EmailVerification) error {
	v.ID = m.nextID
	m.nextID++
	m.verifications = append(m.verifications, *v)
	return nil
}

func (m *mockVerificationRepo) FindLatest(email, purpose string) (*model.EmailVerification, error) {
	for i := len(m.verifications) - 1; i >= 0; i-- {
		v := m.verifications[i]
		if v.Email == email && v.Purpose == purpose {
			return &v, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockVerificationRepo) MarkUsed(id uint) error {
	for i, v := range m.verifications {
		if v.ID == id {
			m.verifications[i].Used = true
			return nil
		}
	}
	return nil
}

func (m *mockVerificationRepo) CreateSendLog(log *model.SendLog) error {
	log.ID = m.nextID
	m.nextID++
	m.sendLogs = append(m.sendLogs, *log)
	return nil
}

func (m *mockVerificationRepo) CountRecentSends(email, purpose string, within time.Duration) (int64, error) {
	since := time.Now().Add(-within)
	count := int64(0)
	for _, l := range m.sendLogs {
		if l.Email == email && l.Purpose == purpose && l.SentAt.After(since) {
			count++
		}
	}
	return count, nil
}

func TestGenerateCode(t *testing.T) {
	repo := &mockVerificationRepo{}
	secCfg := &config.SecurityConfig{CodeLength: 6}
	svc := NewEmailService(repo, &config.SMTPConfig{}, secCfg)

	code := svc.GenerateCode(6)
	if len(code) != 6 {
		t.Errorf("expected code length 6, got %d", len(code))
	}

	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("expected digit, got %c", c)
		}
	}
}

func TestVerifyCode_Correct(t *testing.T) {
	repo := &mockVerificationRepo{}
	secCfg := &config.SecurityConfig{CodeLength: 6, CodeExpiry: 10 * time.Minute}
	svc := NewEmailService(repo, &config.SMTPConfig{}, secCfg)

	v := &model.EmailVerification{
		Email:     "test@example.com",
		Code:      "123456",
		Purpose:   "register",
		Used:      false,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	repo.Create(v)

	err := svc.VerifyCode("test@example.com", "123456", "register")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestVerifyCode_WrongCode(t *testing.T) {
	repo := &mockVerificationRepo{}
	secCfg := &config.SecurityConfig{CodeLength: 6, CodeExpiry: 10 * time.Minute}
	svc := NewEmailService(repo, &config.SMTPConfig{}, secCfg)

	v := &model.EmailVerification{
		Email:     "test@example.com",
		Code:      "123456",
		Purpose:   "register",
		Used:      false,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	repo.Create(v)

	err := svc.VerifyCode("test@example.com", "999999", "register")
	if err == nil {
		t.Fatal("expected error for wrong code, got nil")
	}
}

func TestVerifyCode_Expired(t *testing.T) {
	repo := &mockVerificationRepo{}
	secCfg := &config.SecurityConfig{CodeLength: 6, CodeExpiry: 10 * time.Minute}
	svc := NewEmailService(repo, &config.SMTPConfig{}, secCfg)

	v := &model.EmailVerification{
		Email:     "test@example.com",
		Code:      "123456",
		Purpose:   "register",
		Used:      false,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}
	repo.Create(v)

	err := svc.VerifyCode("test@example.com", "123456", "register")
	if err == nil {
		t.Fatal("expected error for expired code, got nil")
	}
}

func TestVerifyCode_AlreadyUsed(t *testing.T) {
	repo := &mockVerificationRepo{}
	secCfg := &config.SecurityConfig{CodeLength: 6, CodeExpiry: 10 * time.Minute}
	svc := NewEmailService(repo, &config.SMTPConfig{}, secCfg)

	v := &model.EmailVerification{
		Email:     "test@example.com",
		Code:      "123456",
		Purpose:   "register",
		Used:      true,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	repo.Create(v)

	err := svc.VerifyCode("test@example.com", "123456", "register")
	if err == nil {
		t.Fatal("expected error for used code, got nil")
	}
}

func TestRateLimit(t *testing.T) {
	repo := &mockVerificationRepo{}
	secCfg := &config.SecurityConfig{CodeLength: 6, CodeExpiry: 10 * time.Minute, MaxSendPerHour: 3}
	smtpCfg := &config.SMTPConfig{Enabled: true, Host: "smtp.test.com", Port: 587, Username: "test@test.com", Password: "pass"}
	svc := NewEmailService(repo, smtpCfg, secCfg)

	for i := 0; i < 3; i++ {
		repo.CreateSendLog(&model.SendLog{
			Email:   "test@example.com",
			Purpose: "register",
			SentAt:  time.Now(),
		})
	}

	err := svc.SendVerificationCode("test@example.com", "register")
	if err == nil {
		t.Fatal("expected rate limit error, got nil")
	}
}
