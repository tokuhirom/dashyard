package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/prometheus"
)

// HealthHandler handles GET /health - returns server and upstream health status.
type HealthHandler struct {
	client *prometheus.Client
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(client *prometheus.Client) *HealthHandler {
	return &HealthHandler{client: client}
}

// Handle returns the health status of the server and Prometheus connectivity.
func (h *HealthHandler) Handle(c *gin.Context) {
	promStatus := "reachable"
	status := http.StatusOK

	if err := h.client.Ping(c.Request.Context()); err != nil {
		promStatus = "unreachable"
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"status":     statusText(status),
		"prometheus": promStatus,
	})
}

func statusText(code int) string {
	if code == http.StatusOK {
		return "ok"
	}
	return "degraded"
}
