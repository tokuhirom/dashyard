package config

import (
	"testing"
	"time"
)

func TestParseFullConfig(t *testing.T) {
	input := []byte(`
site_title: "My Monitoring"

server:
  host: "127.0.0.1"
  port: 9090
  session_secret: "my-secret"

prometheus:
  url: "http://prom:9090"
  timeout: 60s

dashboards:
  dir: "/etc/dashboards"

users:
  - id: "admin"
    password_hash: "$6$salt$hash"
  - id: "viewer"
    password_hash: "$6$salt2$hash2"
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SiteTitle != "My Monitoring" {
		t.Errorf("expected site_title 'My Monitoring', got %q", cfg.SiteTitle)
	}
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected host '127.0.0.1', got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Server.Port)
	}
	if cfg.Server.SessionSecret != "my-secret" {
		t.Errorf("expected session_secret 'my-secret', got %q", cfg.Server.SessionSecret)
	}
	if cfg.Prometheus.URL != "http://prom:9090" {
		t.Errorf("expected prometheus url 'http://prom:9090', got %q", cfg.Prometheus.URL)
	}
	if cfg.Prometheus.Timeout != 60*time.Second {
		t.Errorf("expected timeout 60s, got %v", cfg.Prometheus.Timeout)
	}
	if cfg.Dashboards.Dir != "/etc/dashboards" {
		t.Errorf("expected dashboards dir '/etc/dashboards', got %q", cfg.Dashboards.Dir)
	}
	if len(cfg.Users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(cfg.Users))
	}
	if cfg.Users[0].ID != "admin" {
		t.Errorf("expected first user 'admin', got %q", cfg.Users[0].ID)
	}
}

func TestParseDefaults(t *testing.T) {
	input := []byte(`{}`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SiteTitle != "Dashyard" {
		t.Errorf("expected default site_title 'Dashyard', got %q", cfg.SiteTitle)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected default host '0.0.0.0', got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Prometheus.URL != "http://localhost:9090" {
		t.Errorf("expected default prometheus url, got %q", cfg.Prometheus.URL)
	}
	if cfg.Prometheus.Timeout != 30*time.Second {
		t.Errorf("expected default timeout 30s, got %v", cfg.Prometheus.Timeout)
	}
	if cfg.Dashboards.Dir != "dashboards" {
		t.Errorf("expected default dashboards dir 'dashboards', got %q", cfg.Dashboards.Dir)
	}
	// Session secret should be auto-generated
	if cfg.Server.SessionSecret == "" {
		t.Error("expected auto-generated session secret")
	}
	if len(cfg.Server.SessionSecret) != 64 { // 32 bytes hex-encoded
		t.Errorf("expected 64-char hex session secret, got %d chars", len(cfg.Server.SessionSecret))
	}
}

func TestParseInvalidYAML(t *testing.T) {
	input := []byte(`{invalid yaml`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}
