package handler

import (
	"zhanbu/internal/service"
	"zhanbu/pkg/response"

	"github.com/gin-gonic/gin"
)

// HoroscopeHandler handles horoscope HTTP requests.
type HoroscopeHandler struct {
	service *service.HoroscopeService
}

// NewHoroscopeHandler creates a new horoscope handler.
func NewHoroscopeHandler(svc *service.HoroscopeService) *HoroscopeHandler {
	return &HoroscopeHandler{service: svc}
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

	response.Success(c, result)
}
