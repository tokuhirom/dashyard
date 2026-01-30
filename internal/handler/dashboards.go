package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/dashboard"
)

// DashboardsHandler handles dashboard listing and detail requests.
type DashboardsHandler struct {
	store     *dashboard.Store
	siteTitle string
}

// NewDashboardsHandler creates a new DashboardsHandler.
func NewDashboardsHandler(store *dashboard.Store, siteTitle string) *DashboardsHandler {
	return &DashboardsHandler{store: store, siteTitle: siteTitle}
}

// List handles GET /api/dashboards - returns all dashboards with flat list and tree.
func (h *DashboardsHandler) List(c *gin.Context) {
	type listItem struct {
		Path  string `json:"path"`
		Title string `json:"title"`
	}

	dashboards := h.store.List()
	items := make([]listItem, len(dashboards))
	for i, d := range dashboards {
		items[i] = listItem{Path: d.Path, Title: d.Title}
	}

	c.JSON(http.StatusOK, gin.H{
		"dashboards": items,
		"tree":       h.store.Tree(),
		"site_title": h.siteTitle,
	})
}

// Get handles GET /api/dashboards/:path - returns a single dashboard definition.
func (h *DashboardsHandler) Get(c *gin.Context) {
	path := c.Param("path")
	// Gin captures the wildcard with a leading slash, strip it
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	d := h.store.Get(path)
	if d == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dashboard not found"})
		return
	}

	c.JSON(http.StatusOK, d)
}

// GetSource handles GET /api/dashboard-source/:path - returns raw YAML source.
func (h *DashboardsHandler) GetSource(c *gin.Context) {
	path := c.Param("path")
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	src, ok := h.store.GetSource(path)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "dashboard not found"})
		return
	}

	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(src))
}
