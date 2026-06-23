package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"zhanbu/internal/middleware"
	"zhanbu/internal/service"
	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/response"
)

// HistoryHandler handles history-related HTTP requests.
type HistoryHandler struct {
	service *service.HistoryService
}

// NewHistoryHandler creates a new HistoryHandler.
func NewHistoryHandler(svc *service.HistoryService) *HistoryHandler {
	return &HistoryHandler{service: svc}
}

// List handles GET /api/history
func (h *HistoryHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	divinationType := c.Query("type")
	if divinationType != "" {
		validTypes := map[string]bool{"tarot": true, "horoscope": true, "liuyao": true, "liuyao_v2": true, "bazi": true, "meihua": true}
		if !validTypes[divinationType] {
			response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid type parameter")
			return
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, appErr := h.service.List(userID, divinationType, page, pageSize)
	if appErr != nil {
		response.ErrorWithAppError(c, http.StatusInternalServerError, appErr)
		return
	}

	response.Success(c, result)
}

// GetDetail handles GET /api/history/:id
func (h *HistoryHandler) GetDetail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid id")
		return
	}

	record, appErr := h.service.GetDetail(userID, uint(id))
	if appErr != nil {
		httpStatus := http.StatusInternalServerError
		if appErr.Code == apperrors.ErrNotFound {
			httpStatus = http.StatusNotFound
		}
		response.ErrorWithAppError(c, httpStatus, appErr)
		return
	}

	response.Success(c, record)
}

// Delete handles DELETE /api/history/:id
func (h *HistoryHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid id")
		return
	}

	appErr := h.service.Delete(userID, uint(id))
	if appErr != nil {
		httpStatus := http.StatusInternalServerError
		if appErr.Code == apperrors.ErrNotFound {
			httpStatus = http.StatusNotFound
		}
		response.ErrorWithAppError(c, httpStatus, appErr)
		return
	}

	response.Success(c, gin.H{"message": "deleted"})
}
