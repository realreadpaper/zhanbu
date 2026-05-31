package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"zhanbu/config"
	"zhanbu/internal/model"
	"zhanbu/internal/repository"
	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/utils"
)

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo    *repository.UserRepository
	jwtManager  *utils.JWTManager
	cfg         *config.JWTConfig
	emailService *EmailService
	secCfg      *config.SecurityConfig
	smtpCfg     *config.SMTPConfig
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	userRepo *repository.UserRepository,
	jwtManager *utils.JWTManager,
	cfg *config.JWTConfig,
	emailService *EmailService,
	secCfg *config.SecurityConfig,
	smtpCfg *config.SMTPConfig,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtManager:  jwtManager,
		cfg:         cfg,
		emailService: emailService,
		secCfg:      secCfg,
		smtpCfg:     smtpCfg,
	}
}

// RegisterRequest represents the registration request body.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents the profile update request body.
type UpdateProfileRequest struct {
	Username  string `json:"username" binding:"omitempty,min=3,max=20"`
	Avatar    string `json:"avatar"`
	Zodiac    string `json:"zodiac"`
	BirthDate string `json:"birth_date"`
}

// RegisterResponse represents the registration response.
type RegisterResponse struct {
	User       model.UserResponse `json:"user"`
	NeedVerify bool               `json:"need_verify"`
}

// Register creates a new user account.
func (s *AuthService) Register(req *RegisterRequest) (*RegisterResponse, *apperrors.AppError) {
	// Validate input
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if err := validateUsername(req.Username); err != nil {
		return nil, err
	}
	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}
	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "database error", err)
	}
	if exists {
		return nil, apperrors.New(apperrors.ErrUserExists, "email already registered")
	}

	exists, err = s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "database error", err)
	}
	if exists {
		return nil, apperrors.New(apperrors.ErrUserExists, "username already taken")
	}

	// Hash password
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "failed to hash password", err)
	}

	// Create user
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "failed to create user", err)
	}

	needVerify := false
	// Send verification email if enabled
	if s.secCfg.VerifyEmail && s.smtpCfg.Enabled {
		if appErr := s.emailService.SendVerificationCode(req.Email, "register"); appErr != nil {
			// Log but don't fail registration — user can resend later
			fmt.Printf("[WARN] failed to send verification email to %s: %v\n", req.Email, appErr)
		}
		needVerify = true
	}

	return &RegisterResponse{
		User:       user.ToResponse(),
		NeedVerify: needVerify,
	}, nil
}

// Login authenticates a user and returns tokens.
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, *apperrors.AppError) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(apperrors.ErrInvalidCreds, "invalid email or password")
		}
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "database error", err)
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, apperrors.New(apperrors.ErrInvalidCreds, "invalid email or password")
	}

	// Check email verification
	if s.secCfg.VerifyEmail && !user.EmailVerified {
		return nil, apperrors.New(apperrors.ErrEmailNotVerified, "请先验证邮箱后再登录")
	}

	// Generate tokens
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "failed to generate tokens", err)
	}

	return &LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         user.ToResponse(),
	}, nil
}

// RefreshToken refreshes an access token using a refresh token.
func (s *AuthService) RefreshToken(refreshTokenStr string) (*LoginResponse, *apperrors.AppError) {
	claims, err := s.jwtManager.ValidateToken(refreshTokenStr)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrTokenInvalid, "invalid or expired refresh token")
	}

	// Verify user still exists
	_, err = s.userRepo.FindByID(claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(apperrors.ErrUnauthorized, "user not found")
		}
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "database error", err)
	}

	// Generate new tokens
	tokens, err := s.jwtManager.GenerateTokenPair(claims.UserID, claims.Username)
	if err != nil {
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "failed to generate tokens", err)
	}

	return &LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User: model.UserResponse{
			ID:       claims.UserID,
			Username: claims.Username,
		},
	}, nil
}

// GetProfile returns the current user's profile.
func (s *AuthService) GetProfile(userID uint) (*model.UserResponse, *apperrors.AppError) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(apperrors.ErrNotFound, "user not found")
		}
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "database error", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateProfile updates the current user's profile.
func (s *AuthService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*model.UserResponse, *apperrors.AppError) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(apperrors.ErrNotFound, "user not found")
		}
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "database error", err)
	}

	fields := make(map[string]interface{})

	if req.Username != "" {
		req.Username = strings.TrimSpace(req.Username)
		if err := validateUsername(req.Username); err != nil {
			return nil, err
		}
		// Check if username is taken by another user
		existing, err := s.userRepo.FindByUsername(req.Username)
		if err == nil && existing.ID != userID {
			return nil, apperrors.New(apperrors.ErrUserExists, "username already taken")
		}
		fields["username"] = req.Username
	}

	if req.Avatar != "" {
		fields["avatar"] = req.Avatar
	}

	if req.Zodiac != "" {
		fields["zodiac"] = req.Zodiac
	}

	if req.BirthDate != "" {
		fields["birth_date"] = req.BirthDate
	}

	if len(fields) > 0 {
		if err := s.userRepo.UpdateFields(user, fields); err != nil {
			return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "failed to update profile", err)
		}
	}

	// Fetch updated user
	updatedUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "database error", err)
	}

	response := updatedUser.ToResponse()
	return &response, nil
}

// VerifyEmailRequest represents the verify email request body.
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// VerifyEmail verifies a user's email with a code.
func (s *AuthService) VerifyEmail(req *VerifyEmailRequest) *apperrors.AppError {
	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Verify the code via email service
	if appErr := s.emailService.VerifyCode(email, req.Code, "register"); appErr != nil {
		return appErr
	}

	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New(apperrors.ErrNotFound, "用户不存在")
		}
		return apperrors.NewWithErr(apperrors.ErrInternalServer, "database error", err)
	}

	// Already verified
	if user.EmailVerified {
		return nil
	}

	// Mark email as verified
	if err := s.userRepo.UpdateFields(user, map[string]interface{}{
		"email_verified": true,
	}); err != nil {
		return apperrors.NewWithErr(apperrors.ErrInternalServer, "更新验证状态失败", err)
	}

	return nil
}

// ResendVerificationRequest represents the resend verification request body.
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResendVerification resends a verification email.
func (s *AuthService) ResendVerification(req *ResendVerificationRequest) *apperrors.AppError {
	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Check if user exists
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New(apperrors.ErrNotFound, "用户不存在")
		}
		return apperrors.NewWithErr(apperrors.ErrInternalServer, "database error", err)
	}

	// Already verified
	if user.EmailVerified {
		return apperrors.New(apperrors.ErrBadRequest, "邮箱已验证")
	}

	// Send verification email
	return s.emailService.SendVerificationCode(email, "register")
}

// LoginResponse represents the login/refresh response.
type LoginResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	ExpiresIn    int64               `json:"expires_in"`
	User         model.UserResponse  `json:"user"`
}

// Validation helpers

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func validateUsername(username string) *apperrors.AppError {
	if len(username) < 3 || len(username) > 20 {
		return apperrors.New(apperrors.ErrInvalidUsername, "username must be 3-20 characters")
	}
	if !usernameRegex.MatchString(username) {
		return apperrors.New(apperrors.ErrInvalidUsername, "username can only contain letters, numbers, underscores and hyphens")
	}
	return nil
}

func validateEmail(email string) *apperrors.AppError {
	if !emailRegex.MatchString(email) {
		return apperrors.New(apperrors.ErrInvalidEmail, "invalid email format")
	}
	return nil
}

func validatePassword(password string) *apperrors.AppError {
	if len(password) < 6 || len(password) > 50 {
		return apperrors.New(apperrors.ErrWeakPassword, "password must be 6-50 characters")
	}
	return nil
}
