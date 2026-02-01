package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/datasource"
)

// ReadyHandler handles GET /ready - checks upstream datasource connectivity.
type ReadyHandler struct {
	registry *datasource.Registry
}

// NewReadyHandler creates a new ReadyHandler.
func NewReadyHandler(registry *datasource.Registry) *ReadyHandler {
	return &ReadyHandler{registry: registry}
}

// Handle returns whether the server is ready to serve traffic.
func (h *ReadyHandler) Handle(c *gin.Context) {
	status := http.StatusOK
	datasources := make(map[string]string)

	for _, name := range h.registry.Names() {
		client, _ := h.registry.Get(name)
		if err := client.Ping(c.Request.Context()); err != nil {
			datasources[name] = "unreachable"
			status = http.StatusServiceUnavailable
		} else {
			datasources[name] = "reachable"
		}
	}

	c.JSON(status, gin.H{
		"status":      readyStatusText(status),
		"datasources": datasources,
	})
}

func readyStatusText(code int) string {
	if code == http.StatusOK {
		return "ok"
	}
	return "degraded"
}
