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

datasources:
  - name: main
    type: prometheus
    url: "http://prom:9090"
    timeout: 60s
    default: true

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
	if len(cfg.Datasources) != 1 {
		t.Fatalf("expected 1 datasource, got %d", len(cfg.Datasources))
	}
	if cfg.Datasources[0].URL != "http://prom:9090" {
		t.Errorf("expected datasource url 'http://prom:9090', got %q", cfg.Datasources[0].URL)
	}
	if cfg.Datasources[0].Timeout != 60*time.Second {
		t.Errorf("expected timeout 60s, got %v", cfg.Datasources[0].Timeout)
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

	// Default datasource should be auto-generated when none configured
	if len(cfg.Datasources) != 1 {
		t.Fatalf("expected 1 default datasource, got %d", len(cfg.Datasources))
	}
	if cfg.Datasources[0].Name != "default" {
		t.Errorf("expected datasource name 'default', got %q", cfg.Datasources[0].Name)
	}
	if cfg.Datasources[0].URL != "http://localhost:9090" {
		t.Errorf("expected datasource url 'http://localhost:9090', got %q", cfg.Datasources[0].URL)
	}
}

func TestParseCookieSecure(t *testing.T) {
	input := []byte(`
server:
  session_secret: "test"
  cookie_secure: true
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.Server.CookieSecure {
		t.Error("expected cookie_secure to be true")
	}
}

func TestParseCookieSecureDefault(t *testing.T) {
	input := []byte(`{}`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.CookieSecure {
		t.Error("expected cookie_secure to default to false")
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

func TestParseDatasources(t *testing.T) {
	input := []byte(`
datasources:
  - name: main
    type: prometheus
    url: "http://prom1:9090"
    timeout: 60s
    default: true
  - name: app
    type: prometheus
    url: "http://prom2:9090"
    timeout: 15s
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Datasources) != 2 {
		t.Fatalf("expected 2 datasources, got %d", len(cfg.Datasources))
	}
	if cfg.Datasources[0].Name != "main" {
		t.Errorf("expected name 'main', got %q", cfg.Datasources[0].Name)
	}
	if cfg.Datasources[0].URL != "http://prom1:9090" {
		t.Errorf("expected url 'http://prom1:9090', got %q", cfg.Datasources[0].URL)
	}
	if cfg.Datasources[0].Timeout != 60*time.Second {
		t.Errorf("expected timeout 60s, got %v", cfg.Datasources[0].Timeout)
	}
	if !cfg.Datasources[0].Default {
		t.Error("expected first datasource to be default")
	}
	if cfg.Datasources[1].Name != "app" {
		t.Errorf("expected name 'app', got %q", cfg.Datasources[1].Name)
	}
	if cfg.Datasources[1].Default {
		t.Error("expected second datasource to not be default")
	}
}

func TestParseSingleDatasourceAutoDefault(t *testing.T) {
	input := []byte(`
datasources:
  - name: solo
    type: prometheus
    url: "http://solo:9090"
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.Datasources[0].Default {
		t.Error("expected single datasource to be auto-set as default")
	}
}

func TestDefaultDatasource(t *testing.T) {
	input := []byte(`
datasources:
  - name: first
    type: prometheus
    url: "http://first:9090"
  - name: second
    type: prometheus
    url: "http://second:9090"
    default: true
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ds := cfg.DefaultDatasource()
	if ds.Name != "second" {
		t.Errorf("expected default datasource 'second', got %q", ds.Name)
	}
}

func TestParseDatasourceValidationDuplicateName(t *testing.T) {
	input := []byte(`
datasources:
  - name: dup
    type: prometheus
    url: "http://a:9090"
    default: true
  - name: dup
    type: prometheus
    url: "http://b:9090"
`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for duplicate datasource name")
	}
}

func TestParseDatasourceValidationUnsupportedType(t *testing.T) {
	input := []byte(`
datasources:
  - name: test
    type: influxdb
    url: "http://a:9090"
`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for unsupported datasource type")
	}
}

func TestParseDatasourceValidationMissingURL(t *testing.T) {
	input := []byte(`
datasources:
  - name: test
    type: prometheus
`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for missing datasource url")
	}
}

func TestParseDatasourceValidationMissingName(t *testing.T) {
	input := []byte(`
datasources:
  - type: prometheus
    url: "http://a:9090"
`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for missing datasource name")
	}
}

func TestParseDatasourceValidationMultipleDefaults(t *testing.T) {
	input := []byte(`
datasources:
  - name: a
    type: prometheus
    url: "http://a:9090"
    default: true
  - name: b
    type: prometheus
    url: "http://b:9090"
    default: true
`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for multiple default datasources")
	}
}

func TestParseDatasourceHeaders(t *testing.T) {
	input := []byte(`
datasources:
  - name: prod
    type: prometheus
    url: "https://prometheus.example.com"
    timeout: 30s
    default: true
    headers:
      - name: Authorization
        value: "Bearer my-secret-token"
      - name: X-Custom-Header
        value: "custom-value"
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Datasources) != 1 {
		t.Fatalf("expected 1 datasource, got %d", len(cfg.Datasources))
	}
	headers := cfg.Datasources[0].Headers
	if len(headers) != 2 {
		t.Fatalf("expected 2 headers, got %d", len(headers))
	}
	if headers[0].Name != "Authorization" || headers[0].Value != "Bearer my-secret-token" {
		t.Errorf("expected Authorization: 'Bearer my-secret-token', got %q: %q", headers[0].Name, headers[0].Value)
	}
	if headers[1].Name != "X-Custom-Header" || headers[1].Value != "custom-value" {
		t.Errorf("expected X-Custom-Header: 'custom-value', got %q: %q", headers[1].Name, headers[1].Value)
	}
}

func TestParseDatasourceHeadersDuplicateKeys(t *testing.T) {
	input := []byte(`
datasources:
  - name: prod
    type: prometheus
    url: "https://prometheus.example.com"
    default: true
    headers:
      - name: X-Custom
        value: "value1"
      - name: X-Custom
        value: "value2"
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	headers := cfg.Datasources[0].Headers
	if len(headers) != 2 {
		t.Fatalf("expected 2 headers (duplicate keys allowed), got %d", len(headers))
	}
	if headers[0].Value != "value1" {
		t.Errorf("expected first value 'value1', got %q", headers[0].Value)
	}
	if headers[1].Value != "value2" {
		t.Errorf("expected second value 'value2', got %q", headers[1].Value)
	}
}

func TestParseDatasourceHeadersEnvExpansion(t *testing.T) {
	t.Setenv("DASHYARD_TEST_TOKEN", "expanded-secret")
	input := []byte(`
datasources:
  - name: prod
    type: prometheus
    url: "https://prometheus.example.com"
    default: true
    headers:
      - name: Authorization
        value: "Bearer ${DASHYARD_TEST_TOKEN}"
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := cfg.Datasources[0].Headers[0].Value
	if got != "Bearer expanded-secret" {
		t.Errorf("expected 'Bearer expanded-secret', got %q", got)
	}
}

func TestParseDatasourceHeadersEnvUnset(t *testing.T) {
	input := []byte(`
datasources:
  - name: prod
    type: prometheus
    url: "https://prometheus.example.com"
    default: true
    headers:
      - name: Authorization
        value: "Bearer ${DASHYARD_UNSET_VAR_12345}"
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := cfg.Datasources[0].Headers[0].Value
	if got != "Bearer " {
		t.Errorf("expected 'Bearer ' (empty expansion), got %q", got)
	}
}

func TestParseDatasourceNoHeaders(t *testing.T) {
	input := []byte(`
datasources:
  - name: local
    type: prometheus
    url: "http://localhost:9090"
    default: true
`)

	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Datasources[0].Headers != nil {
		t.Errorf("expected nil headers, got %v", cfg.Datasources[0].Headers)
	}
}

func TestParseDatasourceValidationNoDefaultMultiple(t *testing.T) {
	input := []byte(`
datasources:
  - name: a
    type: prometheus
    url: "http://a:9090"
  - name: b
    type: prometheus
    url: "http://b:9090"
`)
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error when no default is set with multiple datasources")
	}
}
