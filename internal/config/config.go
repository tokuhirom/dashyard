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

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	SessionSecret string `yaml:"session_secret"`
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
	SiteTitle  string           `yaml:"site_title"`
	Server     ServerConfig     `yaml:"server"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
	Dashboards DashboardsConfig `yaml:"dashboards"`
	Users      []User           `yaml:"users"`
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
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
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

	return cfg, nil
}
