package handler

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/prometheus"
)

// QueryHandler handles GET /api/query - proxies requests to Prometheus.
type QueryHandler struct {
	client *prometheus.Client
}

// NewQueryHandler creates a new QueryHandler.
func NewQueryHandler(client *prometheus.Client) *QueryHandler {
	return &QueryHandler{client: client}
}

// Handle processes a Prometheus query_range proxy request.
func (h *QueryHandler) Handle(c *gin.Context) {
	query := c.Query("query")
	start := c.Query("start")
	end := c.Query("end")
	step := c.Query("step")

	if query == "" || start == "" || end == "" || step == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query, start, end, and step parameters are required"})
		return
	}

	body, statusCode, err := h.client.QueryRange(c.Request.Context(), query, start, end, step)
	if err != nil {
		slog.Error("prometheus query failed", "error", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "prometheus query failed"})
		return
	}
	defer func() { _ = body.Close() }()

	c.Header("Content-Type", "application/json")
	c.Status(statusCode)
	if _, err := io.Copy(c.Writer, body); err != nil {
		slog.Error("failed to stream prometheus response", "error", err)
	}
}
