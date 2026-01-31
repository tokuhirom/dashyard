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

// OAuthProviderConfig holds settings for a single OAuth/OIDC provider.
type OAuthProviderConfig struct {
	Provider     string   `yaml:"provider"`
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	RedirectURL  string   `yaml:"redirect_url"`
	Scopes       []string `yaml:"scopes,omitempty"`
	AllowedUsers []string `yaml:"allowed_users,omitempty"`
	AllowedOrgs  []string `yaml:"allowed_orgs,omitempty"`
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	OAuth []OAuthProviderConfig `yaml:"oauth,omitempty"`
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

func validateOAuthConfig(providers []OAuthProviderConfig) error {
	seen := make(map[string]bool)
	for i, p := range providers {
		if p.Provider == "" {
			return fmt.Errorf("auth.oauth[%d]: provider is required", i)
		}
		if p.ClientID == "" {
			return fmt.Errorf("auth.oauth[%d]: client_id is required", i)
		}
		if p.ClientSecret == "" {
			return fmt.Errorf("auth.oauth[%d]: client_secret is required", i)
		}
		if seen[p.Provider] {
			return fmt.Errorf("auth.oauth[%d]: duplicate provider %q", i, p.Provider)
		}
		seen[p.Provider] = true
	}
	return nil
}
