package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var envBracesRe = regexp.MustCompile(`\$\{([^}]+)\}`)

// expandEnvBraces expands only ${VAR} and ${VAR:-default} syntax, leaving bare $VAR untouched.
// This is safe for values like SHA-512 crypt hashes that contain literal $ characters.
// Returns an error if a referenced environment variable is not set and no default is provided.
func expandEnvBraces(s string) (string, error) {
	var expandErr error
	result := envBracesRe.ReplaceAllStringFunc(s, func(match string) string {
		if expandErr != nil {
			return match
		}
		inner := match[2 : len(match)-1]
		name, defaultVal, hasDefault := strings.Cut(inner, ":-")
		val, ok := os.LookupEnv(name)
		if !ok {
			if hasDefault {
				return defaultVal
			}
			expandErr = fmt.Errorf("environment variable %q is not set", name)
			return match
		}
		return val
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

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
		slog.Warn("server.session_secret is not set; generating a random secret. Sessions will not persist across server restarts.")
		secret := make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			return nil, fmt.Errorf("generating session secret: %w", err)
		}
		cfg.Server.SessionSecret = hex.EncodeToString(secret)
	}

	// Expand environment variables in config values (${VAR} syntax only)
	for i, ds := range cfg.Datasources {
		if v, err := expandEnvBraces(ds.URL); err != nil {
			return nil, fmt.Errorf("datasources[%d].url: %w", i, err)
		} else {
			cfg.Datasources[i].URL = v
		}
		for j, h := range ds.Headers {
			if v, err := expandEnvBraces(h.Value); err != nil {
				return nil, fmt.Errorf("datasources[%d].headers[%d].value: %w", i, j, err)
			} else {
				cfg.Datasources[i].Headers[j].Value = v
			}
		}
	}
	for i, u := range cfg.Users {
		if v, err := expandEnvBraces(u.PasswordHash); err != nil {
			return nil, fmt.Errorf("users[%d].password_hash: %w", i, err)
		} else {
			cfg.Users[i].PasswordHash = v
		}
	}
	for i, p := range cfg.Auth.OAuth {
		if v, err := expandEnvBraces(p.ClientID); err != nil {
			return nil, fmt.Errorf("auth.oauth[%d].client_id: %w", i, err)
		} else {
			cfg.Auth.OAuth[i].ClientID = v
		}
		if v, err := expandEnvBraces(p.ClientSecret); err != nil {
			return nil, fmt.Errorf("auth.oauth[%d].client_secret: %w", i, err)
		} else {
			cfg.Auth.OAuth[i].ClientSecret = v
		}
		if v, err := expandEnvBraces(p.RedirectURL); err != nil {
			return nil, fmt.Errorf("auth.oauth[%d].redirect_url: %w", i, err)
		} else {
			cfg.Auth.OAuth[i].RedirectURL = v
		}
		if v, err := expandEnvBraces(p.BaseURL); err != nil {
			return nil, fmt.Errorf("auth.oauth[%d].base_url: %w", i, err)
		} else {
			cfg.Auth.OAuth[i].BaseURL = v
		}
	}
	if v, err := expandEnvBraces(cfg.Server.SessionSecret); err != nil {
		return nil, fmt.Errorf("server.session_secret: %w", err)
	} else {
		cfg.Server.SessionSecret = v
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
