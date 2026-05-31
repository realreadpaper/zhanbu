package handler

import (
	"net/http"
	"strconv"

	"zhanbu/internal/middleware"
	"zhanbu/internal/model"
	"zhanbu/internal/service"
	"zhanbu/pkg/response"

	"github.com/gin-gonic/gin"
)

// TarotHandler 塔罗牌HTTP处理器
type TarotHandler struct {
	service *service.TarotService
}

// NewTarotHandler 创建塔罗牌处理器
func NewTarotHandler(svc *service.TarotService) *TarotHandler {
	return &TarotHandler{service: svc}
}

// GetCards 获取所有塔罗牌
func (h *TarotHandler) GetCards(c *gin.Context) {
	// 通过service获取repo数据
	spreads := h.service.GetSpreads()
	_ = spreads
	// 直接调用底层repo
	cards, err := h.service.GetAllCards()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5001, "获取牌组失败")
		return
	}

	// 支持分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(cards) {
		start = len(cards)
	}
	if end > len(cards) {
		end = len(cards)
	}

	response.Success(c, gin.H{
		"total":     len(cards),
		"page":      page,
		"page_size": pageSize,
		"items":     cards[start:end],
	})
}

// GetCardByID 获取单张牌
func (h *TarotHandler) GetCardByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 1001, "无效的牌ID")
		return
	}

	card, err := h.service.GetCardByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5001, "获取牌失败")
		return
	}
	if card == nil {
		response.Error(c, http.StatusNotFound, 1004, "牌不存在")
		return
	}

	response.Success(c, card)
}

// GetSpreads 获取牌阵列表
func (h *TarotHandler) GetSpreads(c *gin.Context) {
	spreads := h.service.GetSpreads()
	response.Success(c, spreads)
}

// DrawRequest 抽牌请求
type DrawRequest struct {
	Spread   string `json:"spread" binding:"required"`
	Question string `json:"question"`
}

// Draw 抽牌
func (h *TarotHandler) Draw(c *gin.Context) {
	var req DrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 1001, "请求参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	var result *model.DrawResult
	var err error
	if userID > 0 {
		result, err = h.service.DrawCards(req.Spread, req.Question, userID)
	} else {
		result, err = h.service.DrawCards(req.Spread, req.Question)
	}
	if err != nil {
		response.Error(c, http.StatusBadRequest, 1001, err.Error())
		return
	}

	response.Success(c, result)
}
