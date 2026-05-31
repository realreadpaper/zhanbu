package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"zhanbu/internal/middleware"
	"zhanbu/internal/model"
	"zhanbu/internal/service"
	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/response"
)

// LiuYaoHandler handles liuyao (hexagram) divination HTTP requests.
type LiuYaoHandler struct {
	service *service.LiuYaoService
	saver   service.DivinationRecordSaver
}

// NewLiuYaoHandler creates a new LiuYaoHandler.
func NewLiuYaoHandler(svc *service.LiuYaoService, saver service.DivinationRecordSaver) *LiuYaoHandler {
	return &LiuYaoHandler{service: svc, saver: saver}
}

// ThrowRequest is the request body for liuyao throw.
type ThrowRequest struct {
	Question string `json:"question"`
}

// Throw handles POST /api/liuyao/throw.
func (h *LiuYaoHandler) Throw(c *gin.Context) {
	var req ThrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.service.Throw(req.Question)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, "投掷失败: "+err.Error())
		return
	}

	// Save record if user is authenticated
	userID := middleware.GetUserID(c)
	if userID > 0 && h.saver != nil {
		resultJSON := marshalLiuYaoResult(result)
		record := &model.DivinationRecord{
			UserID:   userID,
			Type:     "liuyao",
			Question: req.Question,
			Result:   resultJSON,
		}
		_ = h.saver.Create(record)
	}

	response.Success(c, result)
}

// GetHexagrams handles GET /api/liuyao/hexagrams.
func (h *LiuYaoHandler) GetHexagrams(c *gin.Context) {
	hexagrams := h.service.GetAllHexagrams()
	briefs := make([]model.HexagramBrief, len(hexagrams))
	for i, hx := range hexagrams {
		briefs[i] = hx.ToBrief()
	}
	response.Success(c, briefs)
}

// marshalLiuYaoResult marshals a LiuYaoResult to JSON string.
func marshalLiuYaoResult(result *model.LiuYaoResult) string {
	data, _ := json.Marshal(result)
	return string(data)
}
