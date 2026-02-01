package handler

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/datasource"
)

// LabelValuesHandler handles GET /api/label-values - proxies label values requests to the datasource.
type LabelValuesHandler struct {
	registry *datasource.Registry
}

// NewLabelValuesHandler creates a new LabelValuesHandler.
func NewLabelValuesHandler(registry *datasource.Registry) *LabelValuesHandler {
	return &LabelValuesHandler{registry: registry}
}

// Handle processes a datasource label values proxy request.
func (h *LabelValuesHandler) Handle(c *gin.Context) {
	label := c.Query("label")
	if label == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "label parameter is required"})
		return
	}

	client, err := h.registry.Get(c.Query("datasource"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match := c.Query("match")

	body, statusCode, err := client.LabelValues(c.Request.Context(), label, match)
	if err != nil {
		slog.Error("datasource label values query failed", "error", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "datasource label values query failed"})
		return
	}
	defer func() { _ = body.Close() }()

	data, err := io.ReadAll(body)
	if err != nil {
		slog.Error("failed to read datasource response", "error", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read datasource response"})
		return
	}

	c.Data(statusCode, "application/json", data)
}
