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
    - provider: github
      client_id: "my-client-id"
      client_secret: "my-client-secret"
      redirect_url: "http://localhost:8080/auth/github/callback"
      scopes: ["user:email", "read:org"]
      allowed_users: ["user1"]
      allowed_orgs: ["my-org"]
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Auth.OAuth) != 1 {
		t.Fatalf("expected 1 oauth provider, got %d", len(cfg.Auth.OAuth))
	}
	p := cfg.Auth.OAuth[0]
	if p.Provider != "github" {
		t.Errorf("expected provider 'github', got %q", p.Provider)
	}
	if p.ClientID != "my-client-id" {
		t.Errorf("expected client_id 'my-client-id', got %q", p.ClientID)
	}
	if p.ClientSecret != "my-client-secret" {
		t.Errorf("expected client_secret 'my-client-secret', got %q", p.ClientSecret)
	}
	if p.RedirectURL != "http://localhost:8080/auth/github/callback" {
		t.Errorf("expected redirect_url, got %q", p.RedirectURL)
	}
	if len(p.Scopes) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(p.Scopes))
	}
	if len(p.AllowedUsers) != 1 || p.AllowedUsers[0] != "user1" {
		t.Errorf("expected allowed_users [user1], got %v", p.AllowedUsers)
	}
	if len(p.AllowedOrgs) != 1 || p.AllowedOrgs[0] != "my-org" {
		t.Errorf("expected allowed_orgs [my-org], got %v", p.AllowedOrgs)
	}
}

func TestParseOAuthConfigWithBaseURL(t *testing.T) {
	input := []byte(`
auth:
  oauth:
    - provider: github
      client_id: "my-client-id"
      client_secret: "my-client-secret"
      redirect_url: "http://localhost:8080/auth/github/callback"
      base_url: "https://ghe.example.com"
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Auth.OAuth) != 1 {
		t.Fatalf("expected 1 oauth provider, got %d", len(cfg.Auth.OAuth))
	}
	p := cfg.Auth.OAuth[0]
	if p.BaseURL != "https://ghe.example.com" {
		t.Errorf("expected base_url 'https://ghe.example.com', got %q", p.BaseURL)
	}
}

func TestParseOAuthConfigWithoutBaseURL(t *testing.T) {
	input := []byte(`
auth:
  oauth:
    - provider: github
      client_id: "my-client-id"
      client_secret: "my-client-secret"
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p := cfg.Auth.OAuth[0]
	if p.BaseURL != "" {
		t.Errorf("expected empty base_url, got %q", p.BaseURL)
	}
}

func TestParseOAuthValidationMissingProvider(t *testing.T) {
	input := []byte(`
auth:
  oauth:
    - client_id: "id"
      client_secret: "secret"
`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for missing provider")
	}
}

func TestParseOAuthValidationMissingClientID(t *testing.T) {
	input := []byte(`
auth:
  oauth:
    - provider: github
      client_secret: "secret"
`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for missing client_id")
	}
}

func TestParseOAuthValidationMissingClientSecret(t *testing.T) {
	input := []byte(`
auth:
  oauth:
    - provider: github
      client_id: "id"
`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for missing client_secret")
	}
}

func TestParseOAuthValidationDuplicateProvider(t *testing.T) {
	input := []byte(`
auth:
  oauth:
    - provider: github
      client_id: "id1"
      client_secret: "secret1"
    - provider: github
      client_id: "id2"
      client_secret: "secret2"
`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for duplicate provider")
	}
}

func TestParseNoOAuthConfig(t *testing.T) {
	input := []byte(`{}`)
	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Auth.OAuth) != 0 {
		t.Errorf("expected no oauth providers, got %d", len(cfg.Auth.OAuth))
	}
}
