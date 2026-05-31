package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"zhanbu/internal/middleware"
	"zhanbu/internal/service"
	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/response"
)

// AIHandler handles AI-related HTTP requests.
type AIHandler struct {
	service   *service.AIService
	rateLimit map[string]*rateBucket
	mu        sync.Mutex
	rate      int
	window    time.Duration
}

// rateBucket tracks request counts for rate limiting.
type rateBucket struct {
	count    int
	resetAt  time.Time
}

// NewAIHandler creates a new AIHandler.
func NewAIHandler(svc *service.AIService) *AIHandler {
	h := &AIHandler{
		service:   svc,
		rateLimit: make(map[string]*rateBucket),
		rate:      5,
		window:    time.Minute,
	}

	// Cleanup expired buckets periodically
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			h.mu.Lock()
			now := time.Now()
			for key, bucket := range h.rateLimit {
				if now.After(bucket.resetAt) {
					delete(h.rateLimit, key)
				}
			}
			h.mu.Unlock()
		}
	}()

	return h
}

// checkRateLimit checks if a request is allowed for the given key.
func (h *AIHandler) checkRateLimit(key string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()
	bucket, exists := h.rateLimit[key]

	if !exists || now.After(bucket.resetAt) {
		h.rateLimit[key] = &rateBucket{
			count:   1,
			resetAt: now.Add(h.window),
		}
		return true
	}

	if bucket.count >= h.rate {
		return false
	}

	bucket.count++
	return true
}

// InterpretRequest is the request body for AI interpretation.
type InterpretRequest struct {
	Type     string `json:"type" binding:"required"`
	ResultID uint   `json:"result_id"`
	Result   string `json:"result"`   // Direct result JSON (alternative to result_id)
	Question string `json:"question"`
}

// Interpret handles POST /api/ai/interpret with SSE streaming.
func (h *AIHandler) Interpret(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	var req InterpretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	// Validate divination type
	validTypes := map[string]bool{
		"tarot": true, "horoscope": true, "liuyao": true, "bazi": true,
	}
	if !validTypes[req.Type] {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid divination type")
		return
	}

	// Rate limiting per user
	rateKey := fmt.Sprintf("ai:%d", userID)
	if !h.checkRateLimit(rateKey) {
		response.Error(c, http.StatusTooManyRequests, apperrors.ErrRateLimited,
			"rate limit exceeded: max 5 requests per minute")
		return
	}

	// Start streaming - either from record or direct result
	var ch <-chan string
	var err error
	if req.ResultID > 0 {
		ch, err = h.service.Interpret(userID, req.ResultID, req.Type, req.Question)
	} else if req.Result != "" {
		ch, err = h.service.InterpretDirect(req.Type, req.Result, req.Question)
	} else {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "either result_id or result is required")
		return
	}
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			httpStatus := http.StatusInternalServerError
			if appErr.Code == apperrors.ErrNotFound {
				httpStatus = http.StatusNotFound
			} else if appErr.Code == apperrors.ErrAIServiceUnavailable {
				httpStatus = http.StatusServiceUnavailable
			}
			response.ErrorWithAppError(c, httpStatus, appErr)
			return
		}
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, "failed to start interpretation")
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// Stream chunks as SSE
	c.Stream(func(w io.Writer) bool {
		chunk, ok := <-ch
		if !ok {
			// Stream ended
			fmt.Fprintf(w, "data: [DONE]\n\n")
			return false
		}

		// Send chunk as SSE
		data, _ := json.Marshal(map[string]string{"text": chunk})
		fmt.Fprintf(w, "data: %s\n\n", data)
		return true
	})
}
