package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"zhanbu/internal/service"
	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/response"
)

// EmailHandler handles email verification requests.
type EmailHandler struct {
	emailService *service.EmailService
}

// NewEmailHandler creates a new EmailHandler.
func NewEmailHandler(emailSvc *service.EmailService) *EmailHandler {
	return &EmailHandler{emailService: emailSvc}
}

// SendCodeRequest is the request body for sending a verification code.
type SendCodeRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Purpose string `json:"purpose" binding:"required"` // register / reset_password
}

// SendCode sends a verification code to the given email.
// POST /api/email/send-code
func (h *EmailHandler) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "参数错误: "+err.Error())
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Validate purpose
	if req.Purpose != "register" && req.Purpose != "reset_password" {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "purpose 必须是 register 或 reset_password")
		return
	}

	appErr := h.emailService.SendVerificationCode(req.Email, req.Purpose)
	if appErr != nil {
		response.ErrorWithAppError(c, http.StatusTooManyRequests, appErr)
		return
	}

	response.Success(c, gin.H{"message": "验证码已发送，请查收邮件"})
}

// VerifyCodeRequest is the request body for verifying a code.
type VerifyCodeRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Code    string `json:"code" binding:"required"`
	Purpose string `json:"purpose" binding:"required"`
}

// VerifyCode verifies the code for an email.
// POST /api/email/verify-code
func (h *EmailHandler) VerifyCode(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "参数错误: "+err.Error())
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	appErr := h.emailService.VerifyCode(req.Email, req.Code, req.Purpose)
	if appErr != nil {
		response.Error(c, http.StatusBadRequest, appErr.Code, appErr.Message)
		return
	}

	response.Success(c, gin.H{"message": "验证成功"})
}
