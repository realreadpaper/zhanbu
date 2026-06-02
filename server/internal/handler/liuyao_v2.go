package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"zhanbu/internal/middleware"
	"zhanbu/internal/model"
	"zhanbu/internal/service"
	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/response"
)

// LiuYaoV2Handler handles LiuYao v2 (Takashima) HTTP requests.
type LiuYaoV2Handler struct {
	service  *service.LiuYaoV2Service
	saver    service.DivinationRecordSaver
}

// NewLiuYaoV2Handler creates a new LiuYaoV2Handler.
func NewLiuYaoV2Handler(svc *service.LiuYaoV2Service, saver service.DivinationRecordSaver) *LiuYaoV2Handler {
	return &LiuYaoV2Handler{
		service: svc,
		saver:   saver,
	}
}

// Throw handles POST /api/liuyao/v2/throw.
func (h *LiuYaoV2Handler) Throw(c *gin.Context) {
	var req model.LiuYaoV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	result, err := h.service.Throw(req.Question, req.Method)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, "divination failed: "+err.Error())
		return
	}

	// Save record if user is authenticated
	if userID := middleware.GetUserID(c); userID > 0 && h.saver != nil {
		resultJSON, _ := json.Marshal(result)
		record := &model.DivinationRecord{
			UserID:   userID,
			Type:     "liuyao_v2",
			Question: req.Question,
			Result:   string(resultJSON),
		}
		_ = h.saver.Create(record)
	}

	response.Success(c, result)
}

// GetHexagrams handles GET /api/liuyao/v2/hexagrams.
func (h *LiuYaoV2Handler) GetHexagrams(c *gin.Context) {
	hexagrams := h.service.GetHexagrams()
	response.Success(c, hexagrams)
}

// GetHexagramByID handles GET /api/liuyao/v2/hexagrams/:id.
func (h *LiuYaoV2Handler) GetHexagramByID(c *gin.Context) {
	id := 0
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid hexagram id")
		return
	}

	hexagram, err := h.service.GetHexagramByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, apperrors.ErrNotFound, "hexagram not found")
		return
	}

	response.Success(c, hexagram)
}

// GetConfig handles GET /api/liuyao/v2/config.
func (h *LiuYaoV2Handler) GetConfig(c *gin.Context) {
	response.Success(c, gin.H{
		"version": "v2",
		"method":  h.service.GetMethod(),
	})
}
