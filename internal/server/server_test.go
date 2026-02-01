package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"

	"github.com/tokuhirom/dashyard/internal/config"
	"github.com/tokuhirom/dashyard/internal/dashboard"
)

func minimalConfig() *config.Config {
	return &config.Config{
		SiteTitle: "Test",
		Server: config.ServerConfig{
			SessionSecret: "test-secret-that-is-at-least-32-bytes-long!",
		},
		Datasources: []config.DatasourceConfig{
			{Name: "default", Type: "prometheus", URL: "http://localhost:9090", Timeout: 5 * time.Second, Default: true},
		},
	}
}

func emptyHolder() *dashboard.StoreHolder {
	store, _ := dashboard.LoadDir("testdata")
	if store == nil {
		// LoadDir fails if dir doesn't exist; create an empty holder manually.
		store = &dashboard.Store{}
	}
	return dashboard.NewStoreHolder(store)
}

func emptyFS() fstest.MapFS {
	return fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
	}
}

func TestNewServerSuccess(t *testing.T) {
	cfg := minimalConfig()
	srv, err := New(cfg, emptyHolder(), emptyFS(), "127.0.0.1", 0, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestNewServerWithMetrics(t *testing.T) {
	cfg := minimalConfig()
	srv, err := New(cfg, emptyHolder(), emptyFS(), "127.0.0.1", 0, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// /metrics endpoint should be registered
	req := httptest.NewRequest("GET", "/metrics", nil)
	resp := httptest.NewRecorder()
	srv.Handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200 for /metrics, got %d", resp.Code)
	}
}

func TestNewServerInvalidDatasourceType(t *testing.T) {
	cfg := minimalConfig()
	cfg.Datasources = []config.DatasourceConfig{
		{Name: "bad", Type: "influxdb", URL: "http://localhost:8086", Timeout: 5 * time.Second, Default: true},
	}

	_, err := New(cfg, emptyHolder(), emptyFS(), "127.0.0.1", 0, false)
	if err == nil {
		t.Fatal("expected error for unsupported datasource type")
	}
}

func TestPublicRoutesWithoutAuth(t *testing.T) {
	cfg := minimalConfig()
	srv, err := New(cfg, emptyHolder(), emptyFS(), "127.0.0.1", 0, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		method string
		path   string
		expect int
	}{
		{"GET", "/ready", http.StatusServiceUnavailable}, // datasource unreachable
		{"GET", "/api/auth-info", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			resp := httptest.NewRecorder()
			srv.Handler.ServeHTTP(resp, req)

			if resp.Code != tt.expect {
				t.Errorf("expected %d, got %d", tt.expect, resp.Code)
			}
		})
	}
}

func TestAuthenticatedRoutesRequireAuth(t *testing.T) {
	cfg := minimalConfig()
	srv, err := New(cfg, emptyHolder(), emptyFS(), "127.0.0.1", 0, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	paths := []string{
		"/api/dashboards",
		"/api/query?query=up&start=1&end=2&step=1s",
		"/api/label-values?label=job",
		"/api/datasources",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest("GET", path, nil)
			resp := httptest.NewRecorder()
			srv.Handler.ServeHTTP(resp, req)

			if resp.Code != http.StatusUnauthorized {
				t.Errorf("expected 401 for %s without auth, got %d", path, resp.Code)
			}
		})
	}
}

func TestServerAddress(t *testing.T) {
	cfg := minimalConfig()
	srv, err := New(cfg, emptyHolder(), emptyFS(), "0.0.0.0", 8080, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if srv.Addr != "0.0.0.0:8080" {
		t.Errorf("expected addr '0.0.0.0:8080', got %q", srv.Addr)
	}
}

func TestSPAFallback(t *testing.T) {
	cfg := minimalConfig()
	srv, err := New(cfg, emptyHolder(), emptyFS(), "127.0.0.1", 0, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Non-existent path should return the SPA index.html
	req := httptest.NewRequest("GET", "/some/unknown/path", nil)
	resp := httptest.NewRecorder()
	srv.Handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200 for SPA fallback, got %d", resp.Code)
	}
}
