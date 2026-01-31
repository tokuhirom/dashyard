package config

import (
	"testing"
	"time"
)

func TestParseFullConfig(t *testing.T) {
	input := []byte(`
site_title: "My Monitoring"
header_color: "#dc2626"

server:
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
	if cfg.HeaderColor != "#dc2626" {
		t.Errorf("expected header_color '#dc2626', got %q", cfg.HeaderColor)
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
	if cfg.HeaderColor != "" {
		t.Errorf("expected default header_color '', got %q", cfg.HeaderColor)
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

func TestParseTrustedProxies(t *testing.T) {
	input := []byte(`
server:
  session_secret: "test"
  trusted_proxies:
    - "10.0.0.1"
    - "10.0.0.2"
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Server.TrustedProxies) != 2 {
		t.Fatalf("expected 2 trusted proxies, got %d", len(cfg.Server.TrustedProxies))
	}
	if cfg.Server.TrustedProxies[0] != "10.0.0.1" {
		t.Errorf("expected first trusted proxy '10.0.0.1', got %q", cfg.Server.TrustedProxies[0])
	}
	if cfg.Server.TrustedProxies[1] != "10.0.0.2" {
		t.Errorf("expected second trusted proxy '10.0.0.2', got %q", cfg.Server.TrustedProxies[1])
	}
}

func TestParseDefaultsNoTrustedProxies(t *testing.T) {
	input := []byte(`{}`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.TrustedProxies != nil {
		t.Errorf("expected nil trusted_proxies, got %v", cfg.Server.TrustedProxies)
	}
}

func TestParseInvalidYAML(t *testing.T) {
	input := []byte(`{invalid yaml`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParseOAuthConfig(t *testing.T) {
	input := []byte(`
auth:
  oauth:
    provider: github
    client_id: "test-id"
    client_secret: "test-secret"
    redirect_url: "http://localhost:8080/auth/callback"
    scopes: ["read:user"]
    allowed_users: ["user1"]
    allowed_orgs: ["my-org"]
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Auth.OAuth == nil {
		t.Fatal("expected oauth config to be set")
	}
	if cfg.Auth.OAuth.Provider != "github" {
		t.Errorf("expected provider 'github', got %q", cfg.Auth.OAuth.Provider)
	}
	if cfg.Auth.OAuth.ClientID != "test-id" {
		t.Errorf("expected client_id 'test-id', got %q", cfg.Auth.OAuth.ClientID)
	}
	if cfg.Auth.OAuth.ClientSecret != "test-secret" {
		t.Errorf("expected client_secret 'test-secret', got %q", cfg.Auth.OAuth.ClientSecret)
	}
	if cfg.Auth.OAuth.RedirectURL != "http://localhost:8080/auth/callback" {
		t.Errorf("expected redirect_url, got %q", cfg.Auth.OAuth.RedirectURL)
	}
	if len(cfg.Auth.OAuth.Scopes) != 1 || cfg.Auth.OAuth.Scopes[0] != "read:user" {
		t.Errorf("expected scopes [read:user], got %v", cfg.Auth.OAuth.Scopes)
	}
	if len(cfg.Auth.OAuth.AllowedUsers) != 1 || cfg.Auth.OAuth.AllowedUsers[0] != "user1" {
		t.Errorf("expected allowed_users [user1], got %v", cfg.Auth.OAuth.AllowedUsers)
	}
	if len(cfg.Auth.OAuth.AllowedOrgs) != 1 || cfg.Auth.OAuth.AllowedOrgs[0] != "my-org" {
		t.Errorf("expected allowed_orgs [my-org], got %v", cfg.Auth.OAuth.AllowedOrgs)
	}
}

func TestParseOAuthConfigOIDC(t *testing.T) {
	input := []byte(`
auth:
  oauth:
    provider: oidc
    client_id: "test-id"
    client_secret: "test-secret"
    issuer_url: "https://accounts.google.com"
    redirect_url: "http://localhost:8080/auth/callback"
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Auth.OAuth.IssuerURL != "https://accounts.google.com" {
		t.Errorf("expected issuer_url, got %q", cfg.Auth.OAuth.IssuerURL)
	}
}

func TestParseOAuthConfigValidation(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "invalid provider",
			input: `
auth:
  oauth:
    provider: invalid
    client_id: "id"
    client_secret: "secret"
    redirect_url: "http://localhost/callback"`,
		},
		{
			name: "missing client_id",
			input: `
auth:
  oauth:
    provider: github
    client_secret: "secret"
    redirect_url: "http://localhost/callback"`,
		},
		{
			name: "missing client_secret",
			input: `
auth:
  oauth:
    provider: github
    client_id: "id"
    redirect_url: "http://localhost/callback"`,
		},
		{
			name: "missing redirect_url",
			input: `
auth:
  oauth:
    provider: github
    client_id: "id"
    client_secret: "secret"`,
		},
		{
			name: "oidc missing issuer_url",
			input: `
auth:
  oauth:
    provider: oidc
    client_id: "id"
    client_secret: "secret"
    redirect_url: "http://localhost/callback"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.input))
			if err == nil {
				t.Error("expected validation error")
			}
		})
	}
}

func TestParseNoOAuthConfig(t *testing.T) {
	input := []byte(`{}`)
	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Auth.OAuth != nil {
		t.Error("expected nil oauth config")
	}
}
