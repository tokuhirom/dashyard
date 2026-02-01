package handler

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/datasource"
)

// QueryHandler handles GET /api/query - proxies requests to the datasource.
type QueryHandler struct {
	registry *datasource.Registry
}

// NewQueryHandler creates a new QueryHandler.
func NewQueryHandler(registry *datasource.Registry) *QueryHandler {
	return &QueryHandler{registry: registry}
}

// Handle processes a datasource query_range proxy request.
func (h *QueryHandler) Handle(c *gin.Context) {
	query := c.Query("query")
	start := c.Query("start")
	end := c.Query("end")
	step := c.Query("step")

	if query == "" || start == "" || end == "" || step == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query, start, end, and step parameters are required"})
		return
	}

	client, err := h.registry.Get(c.Query("datasource"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, statusCode, err := client.QueryRange(c.Request.Context(), query, start, end, step)
	if err != nil {
		slog.Error("datasource query failed", "error", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "datasource query failed"})
		return
	}
	defer func() { _ = body.Close() }()

	c.Header("Content-Type", "application/json")
	c.Status(statusCode)
	if _, err := io.Copy(c.Writer, body); err != nil {
		slog.Error("failed to stream datasource response", "error", err)
	}
}
