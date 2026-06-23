package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"zhanbu/internal/middleware"
	"zhanbu/internal/model"
	"zhanbu/internal/service"
	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/response"
)

// ChatHandler handles chat-related HTTP requests.
type ChatHandler struct {
	service *service.ChatService
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(svc *service.ChatService) *ChatHandler {
	return &ChatHandler{service: svc}
}

// CreateSession handles POST /api/chat/sessions.
func (h *ChatHandler) CreateSession(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	var req struct {
		RecordID uint   `json:"record_id"`
		Type     string `json:"type"`
		Question string `json:"question"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	var session *model.ChatSession
	var err error
	if req.RecordID > 0 {
		session, err = h.service.CreateSession(userID, req.RecordID)
	} else {
		session, err = h.service.CreateModeSession(userID, service.ChatStartOptions{
			Type:     req.Type,
			Question: req.Question,
		})
	}
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			httpStatus := http.StatusInternalServerError
			if appErr.Code == apperrors.ErrNotFound {
				httpStatus = http.StatusNotFound
			} else if appErr.Code == apperrors.ErrAIServiceUnavailable {
				httpStatus = http.StatusServiceUnavailable
			} else if appErr.Code == apperrors.ErrBadRequest {
				httpStatus = http.StatusBadRequest
			}
			response.ErrorWithAppError(c, httpStatus, appErr)
			return
		}
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, "failed to create session")
		return
	}

	response.Success(c, session.ToResponse())
}

// GetSession handles GET /api/chat/sessions/:id.
func (h *ChatHandler) GetSession(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid session id")
		return
	}

	session, err := h.service.GetSession(userID, uint(sessionID))
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.ErrorWithAppError(c, http.StatusNotFound, appErr)
			return
		}
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, "failed to get session")
		return
	}

	response.Success(c, session.ToResponse())
}

// ListSessions handles GET /api/chat/sessions.
func (h *ChatHandler) ListSessions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	sessions, total, err := h.service.ListSessions(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, "failed to list sessions")
		return
	}

	// Convert to response
	sessionResponses := make([]interface{}, len(sessions))
	for i, s := range sessions {
		sessionResponses[i] = s.ToResponse()
	}

	response.Success(c, gin.H{
		"sessions": sessionResponses,
		"total":    total,
		"page":     page,
		"size":     size,
	})
}

// DeleteSession handles DELETE /api/chat/sessions/:id.
func (h *ChatHandler) DeleteSession(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid session id")
		return
	}

	if err := h.service.DeleteSession(userID, uint(sessionID)); err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.ErrorWithAppError(c, http.StatusNotFound, appErr)
			return
		}
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, "failed to delete session")
		return
	}

	response.Success(c, gin.H{"success": true})
}

// SendMessage handles POST /api/chat/sessions/:id/messages with SSE streaming.
func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid session id")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	// Send message and get streaming response
	ch, err := h.service.SendMessage(userID, uint(sessionID), req.Content)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			httpStatus := http.StatusInternalServerError
			if appErr.Code == apperrors.ErrNotFound {
				httpStatus = http.StatusNotFound
			} else if appErr.Code == apperrors.ErrAIServiceUnavailable {
				httpStatus = http.StatusServiceUnavailable
			} else if appErr.Code == apperrors.ErrBadRequest {
				httpStatus = http.StatusBadRequest
			}
			response.ErrorWithAppError(c, httpStatus, appErr)
			return
		}
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, "failed to send message")
		return
	}

	streamSSE(c, ch)
}

// StreamInitialReading handles POST /api/chat/sessions/:id/initial-reading.
func (h *ChatHandler) StreamInitialReading(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized, "unauthorized")
		return
	}

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrBadRequest, "invalid session id")
		return
	}

	ch, err := h.service.StreamInitialReading(userID, uint(sessionID))
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			httpStatus := http.StatusInternalServerError
			if appErr.Code == apperrors.ErrNotFound {
				httpStatus = http.StatusNotFound
			} else if appErr.Code == apperrors.ErrAIServiceUnavailable {
				httpStatus = http.StatusServiceUnavailable
			} else if appErr.Code == apperrors.ErrBadRequest {
				httpStatus = http.StatusBadRequest
			}
			response.ErrorWithAppError(c, httpStatus, appErr)
			return
		}
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternalServer, "failed to stream initial reading")
		return
	}

	streamSSE(c, ch)
}

func streamSSE(c *gin.Context, ch <-chan string) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	c.Stream(func(w io.Writer) bool {
		chunk, ok := <-ch
		if !ok {
			fmt.Fprintf(w, "data: [DONE]\n\n")
			return false
		}

		data, _ := json.Marshal(map[string]string{"text": chunk})
		fmt.Fprintf(w, "data: %s\n\n", data)
		return true
	})
}
