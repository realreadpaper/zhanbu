package service

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net/smtp"
	"strings"
	"time"

	"zhanbu/config"
	"zhanbu/internal/model"
	apperrors "zhanbu/pkg/errors"
)

// VerificationRepoReader defines the interface for verification data access.
type VerificationRepoReader interface {
	Create(v *model.EmailVerification) error
	FindLatest(email, purpose string) (*model.EmailVerification, error)
	MarkUsed(id uint) error
	CreateSendLog(log *model.SendLog) error
	CountRecentSends(email, purpose string, within time.Duration) (int64, error)
}

// EmailService handles email sending and verification.
type EmailService struct {
	verRepo VerificationRepoReader
	cfg     *config.SMTPConfig
	secCfg  *config.SecurityConfig
}

// NewEmailService creates a new EmailService.
func NewEmailService(
	verRepo VerificationRepoReader,
	cfg *config.SMTPConfig,
	secCfg *config.SecurityConfig,
) *EmailService {
	return &EmailService{
		verRepo: verRepo,
		cfg:     cfg,
		secCfg:  secCfg,
	}
}

// GenerateCode generates a random numeric code.
func (s *EmailService) GenerateCode(length int) string {
	if length <= 0 {
		length = 6
	}
	code := ""
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += fmt.Sprintf("%d", n.Int64())
	}
	return code
}

// SendVerificationCode sends a verification code to the email.
func (s *EmailService) SendVerificationCode(email, purpose string) *apperrors.AppError {
	if !s.cfg.Enabled {
		return apperrors.New(apperrors.ErrInternalServer, "邮件服务未启用")
	}

	// Rate limit check
	count, err := s.verRepo.CountRecentSends(email, purpose, 1*time.Hour)
	if err != nil {
		return apperrors.NewWithErr(apperrors.ErrInternalServer, "查询发送记录失败", err)
	}
	if count >= int64(s.secCfg.MaxSendPerHour) {
		return apperrors.New(apperrors.ErrRateLimited, "发送过于频繁，请稍后再试")
	}

	// Generate code
	code := s.GenerateCode(s.secCfg.CodeLength)
	expiresAt := time.Now().Add(s.secCfg.CodeExpiry)

	// Save to database
	verification := &model.EmailVerification{
		Email:     strings.ToLower(strings.TrimSpace(email)),
		Code:      code,
		Purpose:   purpose,
		Used:      false,
		ExpiresAt: expiresAt,
	}
	if err := s.verRepo.Create(verification); err != nil {
		return apperrors.NewWithErr(apperrors.ErrInternalServer, "保存验证码失败", err)
	}

	// Log send
	sendLog := &model.SendLog{
		Email:   email,
		Purpose: purpose,
		SentAt:  time.Now(),
	}
	_ = s.verRepo.CreateSendLog(sendLog)

	// Send email
	subject := "【占卜网】验证码"
	body := fmt.Sprintf(`
<div style="font-family: 'Microsoft YaHei', sans-serif; max-width: 500px; margin: 0 auto; padding: 20px;">
  <h2 style="color: #7c3aed;">🔮 占卜网</h2>
  <p>您的验证码是：</p>
  <div style="background: #f3f4f6; padding: 15px; text-align: center; border-radius: 8px; margin: 20px 0;">
    <span style="font-size: 32px; font-weight: bold; letter-spacing: 8px; color: #7c3aed;">%s</span>
  </div>
  <p style="color: #6b7280;">验证码 %d 分钟内有效，请勿泄露给他人。</p>
  <p style="color: #9ca3af; font-size: 12px;">如非本人操作，请忽略此邮件。</p>
</div>
`, code, int(s.secCfg.CodeExpiry.Minutes()))

	if err := s.sendEmail(email, subject, body); err != nil {
		fmt.Printf("[ERROR] SMTP send failed: host=%s port=%d ssl=%v err=%v\n", s.cfg.Host, s.cfg.Port, s.cfg.SSL, err)
		return apperrors.NewWithErr(apperrors.ErrInternalServer, "发送邮件失败", err)
	}

	return nil
}

// VerifyCode verifies the code for an email.
func (s *EmailService) VerifyCode(email, code, purpose string) *apperrors.AppError {
	email = strings.ToLower(strings.TrimSpace(email))

	verification, err := s.verRepo.FindLatest(email, purpose)
	if err != nil {
		return apperrors.New(apperrors.ErrBadRequest, "验证码无效")
	}

	if verification.Used {
		return apperrors.New(apperrors.ErrBadRequest, "验证码已使用")
	}

	if verification.IsExpired() {
		return apperrors.New(apperrors.ErrBadRequest, "验证码已过期，请重新获取")
	}

	if verification.Code != code {
		return apperrors.New(apperrors.ErrBadRequest, "验证码错误")
	}

	if err := s.verRepo.MarkUsed(verification.ID); err != nil {
		return apperrors.NewWithErr(apperrors.ErrInternalServer, "更新验证码状态失败", err)
	}

	return nil
}

// sendEmail sends an email via SMTP.
func (s *EmailService) sendEmail(to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	from := s.cfg.From
	if from == "" {
		from = s.cfg.Username
	}

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", from, to, subject, htmlBody)

	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	// Port 465 uses implicit TLS (SMTPS), which Go's smtp.SendMail doesn't support.
	// For port 465, establish a TLS connection first, then create an SMTP client.
	// Port 587 uses STARTTLS, which Go's smtp.SendMail handles natively.
	if s.cfg.Port == 465 {
		return s.sendEmailTLS(addr, auth, to, []byte(msg))
	}

	return smtp.SendMail(addr, auth, s.cfg.Username, []string{to}, []byte(msg))
}

// sendEmailTLS sends an email via implicit TLS (SMTPS, port 465).
func (s *EmailService) sendEmailTLS(addr string, auth smtp.Auth, to string, msg []byte) error {
	tlsConfig := &tls.Config{
		ServerName: s.cfg.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS dial failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("SMTP auth failed: %w", err)
			}
		}
	}

	if err := client.Mail(s.cfg.Username); err != nil {
		return fmt.Errorf("SMTP MAIL FROM failed: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT TO failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA failed: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("SMTP write failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("SMTP close failed: %w", err)
	}

	return client.Quit()
}
