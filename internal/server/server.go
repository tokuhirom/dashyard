package server

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/auth"
	"github.com/tokuhirom/dashyard/internal/config"
	"github.com/tokuhirom/dashyard/internal/dashboard"
	"github.com/tokuhirom/dashyard/internal/handler"
	"github.com/tokuhirom/dashyard/internal/prometheus"
)

// New creates and configures an http.Server with all routes and middleware.
func New(cfg *config.Config, store *dashboard.Store, frontendFS fs.FS) *http.Server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Session manager
	sm := auth.NewSessionManager(cfg.Server.SessionSecret, false)

	// Prometheus client
	promClient := prometheus.NewClient(cfg.Prometheus.URL, cfg.Prometheus.Timeout)

	// Handlers
	loginHandler := handler.NewLoginHandler(cfg.Users, sm)
	dashboardsHandler := handler.NewDashboardsHandler(store, cfg.SiteTitle, cfg.HeaderColor)
	queryHandler := handler.NewQueryHandler(promClient)
	staticHandler := handler.NewStaticHandler(frontendFS)

	// Public routes
	r.POST("/api/login", loginHandler.Handle)

	// Authenticated API routes
	api := r.Group("/api")
	api.Use(auth.AuthMiddleware(sm))
	{
		api.GET("/dashboards", dashboardsHandler.List)
		api.GET("/dashboards/*path", dashboardsHandler.Get)
		api.GET("/dashboard-source/*path", dashboardsHandler.GetSource)
		api.GET("/query", queryHandler.Handle)
	}

	// Frontend static files (SPA fallback)
	r.NoRoute(staticHandler.Handle)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
