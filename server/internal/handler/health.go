package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// HealthHandler handles health check requests.
type HealthHandler struct {
	version string
	mode    string
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{version: version}
}

// SetMode sets the server mode (debug/release).
func (h *HealthHandler) SetMode(mode string) {
	h.mode = mode
}

// HealthResponse is the health check response.
type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Uptime    string `json:"uptime,omitempty"`
	GoVersion string `json:"go_version,omitempty"`
	NumCPU    int    `json:"num_cpu,omitempty"`
}

// Check handles GET /api/health.
func (h *HealthHandler) Check(c *gin.Context) {
	uptime := time.Since(startTime)

	resp := HealthResponse{
		Status:  "ok",
		Version: h.version,
		Uptime:  uptime.String(),
	}

	// Only expose system info in non-production mode
	if h.mode != "release" {
		resp.GoVersion = runtime.Version()
		resp.NumCPU = runtime.NumCPU()
	}

	c.JSON(http.StatusOK, resp)
}
