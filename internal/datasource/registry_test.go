package datasource

import (
	"testing"
	"time"

	"github.com/tokuhirom/dashyard/internal/config"
)

func TestNewRegistry(t *testing.T) {
	datasources := []config.DatasourceConfig{
		{Name: "main", Type: "prometheus", URL: "http://main:9090", Timeout: 30 * time.Second, Default: true},
		{Name: "app", Type: "prometheus", URL: "http://app:9090", Timeout: 15 * time.Second},
	}

	reg, err := NewRegistry(datasources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reg.DefaultName() != "main" {
		t.Errorf("expected default name 'main', got %q", reg.DefaultName())
	}
	if reg.Default() == nil {
		t.Fatal("expected non-nil default client")
	}
}

func TestRegistryGet(t *testing.T) {
	datasources := []config.DatasourceConfig{
		{Name: "main", Type: "prometheus", URL: "http://main:9090", Timeout: 30 * time.Second, Default: true},
		{Name: "app", Type: "prometheus", URL: "http://app:9090", Timeout: 15 * time.Second},
	}

	reg, err := NewRegistry(datasources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	client, err := reg.Get("main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client for 'main'")
	}

	client, err = reg.Get("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client for 'app'")
	}

	// Empty name returns default
	client, err = reg.Get("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client != reg.Default() {
		t.Error("expected empty name to return default client")
	}

	// Unknown name returns error
	_, err = reg.Get("nonexistent")
	if err == nil {
		t.Error("expected error for unknown datasource name")
	}
}

func TestRegistryNames(t *testing.T) {
	datasources := []config.DatasourceConfig{
		{Name: "beta", Type: "prometheus", URL: "http://b:9090", Timeout: 30 * time.Second},
		{Name: "alpha", Type: "prometheus", URL: "http://a:9090", Timeout: 30 * time.Second, Default: true},
	}

	reg, err := NewRegistry(datasources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	names := reg.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
	if names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("expected sorted names [alpha, beta], got %v", names)
	}
}

func TestNewRegistryUnsupportedType(t *testing.T) {
	datasources := []config.DatasourceConfig{
		{Name: "main", Type: "influxdb", URL: "http://main:8086", Timeout: 30 * time.Second, Default: true},
	}

	_, err := NewRegistry(datasources)
	if err == nil {
		t.Fatal("expected error for unsupported datasource type")
	}
}
