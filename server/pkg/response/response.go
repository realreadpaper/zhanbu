package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "zhanbu/pkg/errors"
)

// Response is the unified API response format.
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Success sends a success response.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    apperrors.Success,
		Data:    data,
		Message: apperrors.GetMessage(apperrors.Success),
	})
}

// Created sends a 201 Created response.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    apperrors.Success,
		Data:    data,
		Message: apperrors.GetMessage(apperrors.Success),
	})
}

// Error sends an error response with the given HTTP status and error code.
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Data:    nil,
		Message: message,
	})
}

// ErrorWithAppError sends an error response from an AppError.
func ErrorWithAppError(c *gin.Context, httpStatus int, appErr *apperrors.AppError) {
	c.JSON(httpStatus, Response{
		Code:    appErr.Code,
		Data:    nil,
		Message: appErr.Message,
	})
}

// BadRequest sends a 400 error.
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, message)
}

// Unauthorized sends a 401 error.
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, message)
}

// Forbidden sends a 403 error.
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, apperrors.ErrForbidden, message)
}

// NotFound sends a 404 error.
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, apperrors.ErrNotFound, message)
}

// RateLimited sends a 429 error.
func RateLimited(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, apperrors.ErrRateLimited, message)
}

// InternalError sends a 500 error.
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, message)
}
