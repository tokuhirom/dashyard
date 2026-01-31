package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// User represents an authenticated user in the config.
type User struct {
	ID           string `yaml:"id"`
	PasswordHash string `yaml:"password_hash"`
}

// OAuthConfig holds OAuth/OIDC provider settings.
type OAuthConfig struct {
	Provider     string   `yaml:"provider"`      // "github", "google", or "oidc"
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	IssuerURL    string   `yaml:"issuer_url"`    // Required for "oidc", auto-set for "google"
	RedirectURL  string   `yaml:"redirect_url"`
	Scopes       []string `yaml:"scopes"`        // Defaults per provider if omitted
	AllowedUsers []string `yaml:"allowed_users"` // Optional allowlist
	AllowedOrgs  []string `yaml:"allowed_orgs"`  // Optional (GitHub-specific)
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	OAuth *OAuthConfig `yaml:"oauth,omitempty"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	SessionSecret  string   `yaml:"session_secret"`
	TrustedProxies []string `yaml:"trusted_proxies,omitempty"`
}

// PrometheusConfig holds Prometheus connection settings.
type PrometheusConfig struct {
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
}

// DashboardsConfig holds dashboard directory settings.
type DashboardsConfig struct {
	Dir string `yaml:"dir"`
}

// Config is the top-level application configuration.
type Config struct {
	SiteTitle   string           `yaml:"site_title"`
	HeaderColor string           `yaml:"header_color"`
	Server      ServerConfig     `yaml:"server"`
	Prometheus  PrometheusConfig `yaml:"prometheus"`
	Dashboards  DashboardsConfig `yaml:"dashboards"`
	Users       []User           `yaml:"users"`
	Auth        AuthConfig       `yaml:"auth"`
}

// Load reads and parses a YAML config file, applying defaults for missing values.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	return Parse(data)
}

// Parse parses YAML config data, applying defaults for missing values.
func Parse(data []byte) (*Config, error) {
	cfg := &Config{
		SiteTitle: "Dashyard",
		Prometheus: PrometheusConfig{
			URL:     "http://localhost:9090",
			Timeout: 30 * time.Second,
		},
		Dashboards: DashboardsConfig{
			Dir: "dashboards",
		},
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.Server.SessionSecret == "" {
		secret := make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			return nil, fmt.Errorf("generating session secret: %w", err)
		}
		cfg.Server.SessionSecret = hex.EncodeToString(secret)
	}

	if err := validateOAuthConfig(cfg.Auth.OAuth); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateOAuthConfig(oauth *OAuthConfig) error {
	if oauth == nil {
		return nil
	}

	switch oauth.Provider {
	case "github", "google", "oidc":
		// valid
	default:
		return fmt.Errorf("unsupported oauth provider: %q (must be github, google, or oidc)", oauth.Provider)
	}

	if oauth.ClientID == "" {
		return fmt.Errorf("auth.oauth.client_id is required")
	}
	if oauth.ClientSecret == "" {
		return fmt.Errorf("auth.oauth.client_secret is required")
	}
	if oauth.RedirectURL == "" {
		return fmt.Errorf("auth.oauth.redirect_url is required")
	}

	if oauth.Provider == "oidc" && oauth.IssuerURL == "" {
		return fmt.Errorf("auth.oauth.issuer_url is required for oidc provider")
	}

	return nil
}
