package handler

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/prometheus"
)

// LabelValuesHandler handles GET /api/label-values - proxies label values requests to Prometheus.
type LabelValuesHandler struct {
	client *prometheus.Client
}

// NewLabelValuesHandler creates a new LabelValuesHandler.
func NewLabelValuesHandler(client *prometheus.Client) *LabelValuesHandler {
	return &LabelValuesHandler{client: client}
}

// Handle processes a Prometheus label values proxy request.
func (h *LabelValuesHandler) Handle(c *gin.Context) {
	label := c.Query("label")
	if label == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "label parameter is required"})
		return
	}

	match := c.Query("match")

	body, statusCode, err := h.client.LabelValues(c.Request.Context(), label, match)
	if err != nil {
		slog.Error("prometheus label values query failed", "error", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "prometheus label values query failed"})
		return
	}
	defer func() { _ = body.Close() }()

	c.Header("Content-Type", "application/json")
	c.Status(statusCode)
	if _, err := io.Copy(c.Writer, body); err != nil {
		slog.Error("failed to stream prometheus response", "error", err)
	}
}
