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

// BaZiHandler handles bazi (八字) divination HTTP requests.
type BaZiHandler struct {
	service *service.BaZiService
	saver   service.DivinationRecordSaver
}

// NewBaZiHandler creates a new BaZiHandler.
func NewBaZiHandler(svc *service.BaZiService, saver service.DivinationRecordSaver) *BaZiHandler {
	return &BaZiHandler{service: svc, saver: saver}
}

// Calculate handles POST /api/bazi/calculate.
func (h *BaZiHandler) Calculate(c *gin.Context) {
	var req model.BaZiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if req.BirthDate == "" {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "出生日期不能为空")
		return
	}
	if req.BirthTime == "" {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "出生时间不能为空")
		return
	}

	result, err := h.service.Calculate(req.BirthDate, req.BirthTime, req.Gender)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, err.Error())
		return
	}

	// Save record if user is authenticated
	userID := middleware.GetUserID(c)
	if userID > 0 && h.saver != nil {
		resultJSON := marshalBaZiResult(result)
		record := &model.DivinationRecord{
			UserID:   userID,
			Type:     "bazi",
			Question: req.BirthDate + " " + req.BirthTime,
			Result:   resultJSON,
		}
		if err := h.saver.Create(record); err == nil {
			result.RecordID = record.ID
		}
	}

	response.Success(c, result)
}

// marshalBaZiResult marshals a BaZiResult to JSON string.
func marshalBaZiResult(result *model.BaZiResult) string {
	data, _ := json.Marshal(result)
	return string(data)
}
