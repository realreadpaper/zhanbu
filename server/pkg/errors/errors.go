package errors

// Business error codes for the application.
const (
	// Success
	Success = 0

	// General errors (1000-1999)
	ErrBadRequest     = 1001
	ErrUnauthorized   = 1002
	ErrForbidden      = 1003
	ErrNotFound       = 1004
	ErrRateLimited    = 1005
	ErrInternalServer = 1006

	// Auth errors (2000-2999)
	ErrUserExists       = 2001
	ErrInvalidCreds     = 2002
	ErrTokenExpired     = 2003
	ErrTokenInvalid     = 2004
	ErrInvalidUsername  = 2005
	ErrInvalidEmail     = 2006
	ErrWeakPassword      = 2007
	ErrEmailNotVerified  = 2008

	// AI errors (3000-3999)
	ErrAIServiceUnavailable = 3001
	ErrAITimeout           = 3002
)

// CodeMessages maps error codes to default messages.
var CodeMessages = map[int]string{
	Success:               "success",
	ErrBadRequest:         "bad request",
	ErrUnauthorized:       "unauthorized",
	ErrForbidden:          "forbidden",
	ErrNotFound:           "resource not found",
	ErrRateLimited:        "rate limit exceeded",
	ErrInternalServer:     "internal server error",
	ErrUserExists:         "user already exists",
	ErrInvalidCreds:       "invalid email or password",
	ErrTokenExpired:       "token has expired",
	ErrTokenInvalid:       "token is invalid",
	ErrInvalidUsername:    "invalid username format",
	ErrInvalidEmail:       "invalid email format",
	ErrWeakPassword:       "password is too weak",
	ErrEmailNotVerified:   "email not verified",
	ErrAIServiceUnavailable: "AI service is currently unavailable",
	ErrAITimeout:          "AI service request timed out",
}

// AppError is a business-level error with a code and message.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // underlying error, not serialized
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError.
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// NewWithErr creates a new AppError with an underlying error.
func NewWithErr(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// GetMessage returns the default message for a code, or a fallback.
func GetMessage(code int) string {
	if msg, ok := CodeMessages[code]; ok {
		return msg
	}
	return "unknown error"
}
