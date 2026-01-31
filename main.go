package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
	"github.com/alecthomas/kong"
	"github.com/tokuhirom/dashyard/internal/config"
	"github.com/tokuhirom/dashyard/internal/dashboard"
	"github.com/tokuhirom/dashyard/internal/metrics"
	"github.com/tokuhirom/dashyard/internal/server"
)

//go:embed frontend/dist/*
var frontendFiles embed.FS

var cli struct {
	Serve      ServeCmd      `cmd:"" help:"Start the dashboard server."`
	Validate   ValidateCmd   `cmd:"" help:"Validate config and dashboard files."`
	Mkpasswd   MkpasswdCmd   `cmd:"" help:"Generate a SHA-512 crypt password hash."`
	GenPrompt GenPromptCmd `cmd:"gen-prompt" help:"Generate an LLM prompt for dashboard YAML generation from Prometheus metrics."`
}

type ServeCmd struct {
	Config  string `help:"Path to config file." default:"config.yaml"`
	Host    string `help:"Host to listen on." default:"0.0.0.0"`
	Port    int    `help:"Port to listen on." default:"8080"`
	Metrics bool   `help:"Enable /metrics endpoint exposing Prometheus metrics." default:"false"`
}

func (cmd *ServeCmd) Run() error {
	// Load config
	cfg, err := config.Load(cmd.Config)
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
	metrics.DashboardsLoaded.Set(float64(len(store.List())))

	holder := dashboard.NewStoreHolder(store)

	// Get frontend filesystem
	frontendFS, err := fs.Sub(frontendFiles, "frontend/dist")
	if err != nil {
		slog.Error("failed to access frontend files", "error", err)
		os.Exit(1)
	}

	// Create server
	srv, err := server.New(cfg, holder, frontendFS, cmd.Host, cmd.Port, cmd.Metrics)
	if err != nil {
		slog.Error("failed to create server", "error", err)
		os.Exit(1)
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Watch dashboards directory for changes
	watcher := dashboard.NewWatcher(cfg.Dashboards.Dir, holder)
	go func() {
		if err := watcher.Watch(ctx); err != nil {
			slog.Error("dashboard watcher error", "error", err)
		}
	}()

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
	return nil
}

type ValidateCmd struct {
	Config string `help:"Path to config file." default:"config.yaml"`
}

func (cmd *ValidateCmd) Run() error {
	cfg, err := config.Load(cmd.Config)
	if err != nil {
		return fmt.Errorf("config %s: %w", cmd.Config, err)
	}
	fmt.Printf("Config OK: %s\n", cmd.Config)

	store, err := dashboard.LoadDir(cfg.Dashboards.Dir)
	if err != nil {
		return fmt.Errorf("dashboards: %w", err)
	}
	fmt.Printf("Dashboards OK: loaded %d dashboards from %q\n", len(store.List()), cfg.Dashboards.Dir)

	return nil
}

type MkpasswdCmd struct {
	Password string `arg:"" help:"Password to hash."`
}

func (cmd *MkpasswdCmd) Run() error {
	c := crypt.SHA512.New()
	hash, err := c.Generate([]byte(cmd.Password), nil)
	if err != nil {
		return fmt.Errorf("failed to generate hash: %w", err)
	}
	fmt.Println(hash)
	return nil
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("dashyard"),
		kong.Description("Lightweight Prometheus metrics dashboard."),
	)
	if err := ctx.Run(); err != nil {
		slog.Error("error", "error", err)
		os.Exit(1)
	}
}
