package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"zhanbu/internal/middleware"
	"zhanbu/internal/service"
	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/response"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	user, appErr := h.authService.Register(&req)
	if appErr != nil {
		httpStatus := http.StatusInternalServerError
		switch appErr.Code {
		case apperrors.ErrBadRequest, apperrors.ErrInvalidUsername, apperrors.ErrInvalidEmail, apperrors.ErrWeakPassword:
			httpStatus = http.StatusBadRequest
		case apperrors.ErrUserExists:
			httpStatus = http.StatusConflict
		}
		response.ErrorWithAppError(c, httpStatus, appErr)
		return
	}

	response.Created(c, user)
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	result, appErr := h.authService.Login(&req)
	if appErr != nil {
		httpStatus := http.StatusInternalServerError
		switch appErr.Code {
		case apperrors.ErrBadRequest:
			httpStatus = http.StatusBadRequest
		case apperrors.ErrInvalidCreds:
			httpStatus = http.StatusUnauthorized
		case apperrors.ErrEmailNotVerified:
			httpStatus = http.StatusForbidden
		}
		response.ErrorWithAppError(c, httpStatus, appErr)
		return
	}

	response.Success(c, result)
}

// RefreshToken handles POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "refresh_token is required")
		return
	}

	result, appErr := h.authService.RefreshToken(req.RefreshToken)
	if appErr != nil {
		httpStatus := http.StatusInternalServerError
		switch appErr.Code {
		case apperrors.ErrTokenInvalid:
			httpStatus = http.StatusUnauthorized
		case apperrors.ErrUnauthorized:
			httpStatus = http.StatusUnauthorized
		}
		response.ErrorWithAppError(c, httpStatus, appErr)
		return
	}

	response.Success(c, result)
}

// GetProfile handles GET /api/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, appErr := h.authService.GetProfile(userID)
	if appErr != nil {
		httpStatus := http.StatusInternalServerError
		switch appErr.Code {
		case apperrors.ErrNotFound:
			httpStatus = http.StatusNotFound
		}
		response.ErrorWithAppError(c, httpStatus, appErr)
		return
	}

	response.Success(c, user)
}

// UpdateProfile handles PUT /api/auth/profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	user, appErr := h.authService.UpdateProfile(userID, &req)
	if appErr != nil {
		httpStatus := http.StatusInternalServerError
		switch appErr.Code {
		case apperrors.ErrBadRequest, apperrors.ErrInvalidUsername:
			httpStatus = http.StatusBadRequest
		case apperrors.ErrNotFound:
			httpStatus = http.StatusNotFound
		case apperrors.ErrUserExists:
			httpStatus = http.StatusConflict
		}
		response.ErrorWithAppError(c, httpStatus, appErr)
		return
	}

	response.Success(c, user)
}

// VerifyEmail handles POST /api/auth/verify-email
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req service.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	appErr := h.authService.VerifyEmail(&req)
	if appErr != nil {
		httpStatus := http.StatusInternalServerError
		switch appErr.Code {
		case apperrors.ErrBadRequest:
			httpStatus = http.StatusBadRequest
		case apperrors.ErrNotFound:
			httpStatus = http.StatusNotFound
		}
		response.ErrorWithAppError(c, httpStatus, appErr)
		return
	}

	response.Success(c, nil)
}

// ResendVerification handles POST /api/auth/resend-verification
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req service.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	appErr := h.authService.ResendVerification(&req)
	if appErr != nil {
		httpStatus := http.StatusInternalServerError
		switch appErr.Code {
		case apperrors.ErrBadRequest:
			httpStatus = http.StatusBadRequest
		case apperrors.ErrNotFound:
			httpStatus = http.StatusNotFound
		case apperrors.ErrRateLimited:
			httpStatus = http.StatusTooManyRequests
		}
		response.ErrorWithAppError(c, httpStatus, appErr)
		return
	}

	response.Success(c, nil)
}
