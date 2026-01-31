package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/prometheus"
)

// ReadyHandler handles GET /ready - checks upstream Prometheus connectivity.
type ReadyHandler struct {
	client *prometheus.Client
}

// NewReadyHandler creates a new ReadyHandler.
func NewReadyHandler(client *prometheus.Client) *ReadyHandler {
	return &ReadyHandler{client: client}
}

// Handle returns whether the server is ready to serve traffic.
func (h *ReadyHandler) Handle(c *gin.Context) {
	promStatus := "reachable"
	status := http.StatusOK

	if err := h.client.Ping(c.Request.Context()); err != nil {
		promStatus = "unreachable"
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"status":     readyStatusText(status),
		"prometheus": promStatus,
	})
}

func readyStatusText(code int) string {
	if code == http.StatusOK {
		return "ok"
	}
	return "degraded"
}
