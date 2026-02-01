package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/datasource"
)

// DatasourcesHandler handles GET /api/datasources - returns available datasource names.
type DatasourcesHandler struct {
	registry *datasource.Registry
}

// NewDatasourcesHandler creates a new DatasourcesHandler.
func NewDatasourcesHandler(registry *datasource.Registry) *DatasourcesHandler {
	return &DatasourcesHandler{registry: registry}
}

// Handle returns the list of datasource names and the default.
func (h *DatasourcesHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"datasources": h.registry.Names(),
		"default":     h.registry.DefaultName(),
	})
}
