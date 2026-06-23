package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"zhanbu/internal/middleware"
	"zhanbu/internal/model"
	"zhanbu/internal/service"
	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/response"
)

// MeiHuaHandler handles MeiHua divination HTTP requests.
type MeiHuaHandler struct {
	svc   *service.MeiHuaService
	saver service.DivinationRecordSaver
}

// NewMeiHuaHandler creates a new MeiHuaHandler.
func NewMeiHuaHandler(svc *service.MeiHuaService, saver service.DivinationRecordSaver) *MeiHuaHandler {
	return &MeiHuaHandler{svc: svc, saver: saver}
}

// DivinationType 占卜类型标识。
const DivinationType = "meihua"

// DefaultTimezone 默认时区。
const DefaultTimezone = "Asia/Shanghai"

// CalculateByTime handles POST /api/meihua/calculate/time
func (h *MeiHuaHandler) CalculateByTime(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	var req model.MeiHuaTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid request")
		return
	}

	if req.Question == "" {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "question is required")
		return
	}

	tz := req.Timezone
	if tz == "" {
		tz = DefaultTimezone
	}

	result, err := h.svc.CalculateByTime(req.Question, time.Now(), tz)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, err.Error())
		return
	}

	// Save record
	record, saveErr := h.saveRecord(userID, DivinationType, req.Question, result)
	if saveErr != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, "failed to save record")
		return
	}
	result.RecordID = record.ID

	response.Success(c, result)
}

// CalculateByNumbers handles POST /api/meihua/calculate/numbers
func (h *MeiHuaHandler) CalculateByNumbers(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	var req model.MeiHuaNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid request")
		return
	}

	if req.Question == "" {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "question is required")
		return
	}

	if len(req.Numbers) < service.MinNumbersRequired {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest,
			fmt.Sprintf("at least %d numbers required", service.MinNumbersRequired))
		return
	}

	result, err := h.svc.CalculateByNumbers(req.Question, req.Numbers)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, err.Error())
		return
	}

	// Save record
	record, saveErr := h.saveRecord(userID, DivinationType, req.Question, result)
	if saveErr != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, "failed to save record")
		return
	}
	result.RecordID = record.ID

	response.Success(c, result)
}

func (h *MeiHuaHandler) saveRecord(userID uint, divinationType, question string, result any) (*model.DivinationRecord, error) {
	if h.saver == nil {
		return nil, nil
	}
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	record := &model.DivinationRecord{
		UserID:   userID,
		Type:     divinationType,
		Question: question,
		Result:   string(data),
	}
	if err := h.saver.Create(record); err != nil {
		return nil, err
	}
	return record, nil
}
