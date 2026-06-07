package handler

import (
	"encoding/json"

	"zhanbu/internal/middleware"
	"zhanbu/internal/model"
	"zhanbu/internal/service"
	"zhanbu/pkg/response"

	"github.com/gin-gonic/gin"
)

// HoroscopeHandler handles horoscope HTTP requests.
type HoroscopeHandler struct {
	service *service.HoroscopeService
	saver   service.DivinationRecordSaver
}

// NewHoroscopeHandler creates a new horoscope handler.
func NewHoroscopeHandler(svc *service.HoroscopeService, saver ...service.DivinationRecordSaver) *HoroscopeHandler {
	h := &HoroscopeHandler{service: svc}
	if len(saver) > 0 {
		h.saver = saver[0]
	}
	return h
}

// GetHoroscope handles GET /api/horoscope/:zodiac
// Query params: period (daily|weekly|monthly, default daily), date (YYYY-MM-DD, optional).
func (h *HoroscopeHandler) GetHoroscope(c *gin.Context) {
	zodiac := c.Param("zodiac")
	if zodiac == "" {
		response.BadRequest(c, "请指定星座")
		return
	}

	period := c.DefaultQuery("period", "daily")
	date := c.DefaultQuery("date", "")

	result, err := h.service.Generate(zodiac, period, date)
	if err != nil {
		// Determine if it's a bad request or internal error
		response.BadRequest(c, err.Error())
		return
	}

	if userID := middleware.GetUserID(c); userID > 0 && h.saver != nil {
		resultJSON, _ := json.Marshal(result)
		record := &model.DivinationRecord{
			UserID:   userID,
			Type:     "horoscope",
			Question: result.ZodiacCN + " " + result.Period + " " + result.Date,
			Result:   string(resultJSON),
		}
		if err := h.saver.Create(record); err == nil {
			result.RecordID = record.ID
		}
	}

	response.Success(c, result)
}
