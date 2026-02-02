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
	BaseURL      string   `yaml:"base_url,omitempty"`
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
	CookieSecure   bool     `yaml:"cookie_secure"`
	TrustedProxies []string `yaml:"trusted_proxies,omitempty"`
}

// HeaderConfig represents a single HTTP header as a name/value pair.
type HeaderConfig struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// DatasourceConfig holds settings for a single named datasource.
type DatasourceConfig struct {
	Name    string         `yaml:"name"`
	Type    string         `yaml:"type"`
	URL     string         `yaml:"url"`
	Timeout time.Duration  `yaml:"timeout"`
	Default bool           `yaml:"default"`
	Headers []HeaderConfig `yaml:"headers,omitempty"`
}

// DashboardsConfig holds dashboard directory settings.
type DashboardsConfig struct {
	Dir string `yaml:"dir"`
}

// Config is the top-level application configuration.
type Config struct {
	SiteTitle   string             `yaml:"site_title"`
	HeaderColor string             `yaml:"header_color"`
	Server      ServerConfig       `yaml:"server"`
	Datasources []DatasourceConfig `yaml:"datasources"`
	Dashboards  DashboardsConfig   `yaml:"dashboards"`
	Users       []User             `yaml:"users"`
	Auth        AuthConfig         `yaml:"auth"`
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

	// Provide a default datasource when none configured
	if len(cfg.Datasources) == 0 {
		cfg.Datasources = []DatasourceConfig{
			{
				Name:    "default",
				Type:    "prometheus",
				URL:     "http://localhost:9090",
				Timeout: 30 * time.Second,
				Default: true,
			},
		}
	}

	if err := validateDatasources(cfg.Datasources); err != nil {
		return nil, err
	}

	// Expand environment variables in datasource header values
	for i, ds := range cfg.Datasources {
		for j, h := range ds.Headers {
			cfg.Datasources[i].Headers[j].Value = os.ExpandEnv(h.Value)
		}
	}

	return cfg, nil
}

// DefaultDatasource returns the datasource marked as default.
func (c *Config) DefaultDatasource() DatasourceConfig {
	for _, ds := range c.Datasources {
		if ds.Default {
			return ds
		}
	}
	return c.Datasources[0]
}

func validateDatasources(datasources []DatasourceConfig) error {
	if len(datasources) == 0 {
		return fmt.Errorf("at least one datasource must be configured")
	}

	seen := make(map[string]bool)
	defaultCount := 0

	for i, ds := range datasources {
		if ds.Name == "" {
			return fmt.Errorf("datasources[%d]: name is required", i)
		}
		if seen[ds.Name] {
			return fmt.Errorf("datasources[%d]: duplicate name %q", i, ds.Name)
		}
		seen[ds.Name] = true

		validTypes := map[string]bool{"prometheus": true}
		if !validTypes[ds.Type] {
			return fmt.Errorf("datasources[%d]: unsupported type %q", i, ds.Type)
		}
		if ds.URL == "" {
			return fmt.Errorf("datasources[%d]: url is required", i)
		}
		if ds.Default {
			defaultCount++
		}
	}

	// If only one datasource, auto-set as default
	if defaultCount == 0 && len(datasources) == 1 {
		datasources[0].Default = true
		defaultCount = 1
	}

	if defaultCount == 0 {
		return fmt.Errorf("exactly one datasource must be marked as default")
	}
	if defaultCount > 1 {
		return fmt.Errorf("only one datasource can be marked as default, found %d", defaultCount)
	}

	return nil
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
