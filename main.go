package main

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/tokuhirom/dashyard/internal/config"
	"github.com/tokuhirom/dashyard/internal/dashboard"
	"github.com/tokuhirom/dashyard/internal/server"
)

//go:embed frontend/dist/*
var frontendFiles embed.FS

var cli struct {
	Config string `help:"Path to config file." default:"config.yaml"`
}

func main() {
	kong.Parse(&cli,
		kong.Name("dashyard"),
		kong.Description("Lightweight Prometheus metrics dashboard."),
	)

	// Load config
	cfg, err := config.Load(cli.Config)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Load dashboards
	store, err := dashboard.LoadDir(cfg.Dashboards.Dir)
	if err != nil {
		slog.Error("failed to load dashboards", "error", err)
		os.Exit(1)
	}
	slog.Info("loaded dashboards", "count", len(store.List()))

	// Get frontend filesystem
	frontendFS, err := fs.Sub(frontendFiles, "frontend/dist")
	if err != nil {
		slog.Error("failed to access frontend files", "error", err)
		os.Exit(1)
	}

	// Create server
	srv := server.New(cfg, store, frontendFS)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Prometheus.Timeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
	slog.Info("server stopped")
}
