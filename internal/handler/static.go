package handler

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StaticHandler serves embedded frontend files with SPA fallback.
type StaticHandler struct {
	fileServer http.Handler
	fsys       fs.FS
}

// NewStaticHandler creates a new StaticHandler from an embedded filesystem.
func NewStaticHandler(fsys fs.FS) *StaticHandler {
	return &StaticHandler{
		fileServer: http.FileServer(http.FS(fsys)),
		fsys:       fsys,
	}
}

// Handle serves static files, falling back to index.html for SPA routing.
func (h *StaticHandler) Handle(c *gin.Context) {
	path := c.Request.URL.Path

	// Try to serve the exact file
	if f, err := h.fsys.Open(path[1:]); err == nil {
		f.Close()
		h.fileServer.ServeHTTP(c.Writer, c.Request)
		return
	}

	// SPA fallback: serve index.html
	c.Request.URL.Path = "/"
	h.fileServer.ServeHTTP(c.Writer, c.Request)
}
