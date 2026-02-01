package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/config"
	"github.com/tokuhirom/dashyard/internal/datasource"
)

func TestDatasourcesHandler(t *testing.T) {
	registry := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "prod", Type: "prometheus", URL: "http://prom-prod:9090", Timeout: 5 * time.Second, Default: true},
		{Name: "staging", Type: "prometheus", URL: "http://prom-staging:9090", Timeout: 5 * time.Second},
	})
	handler := NewDatasourcesHandler(registry)

	router := gin.New()
	router.GET("/api/datasources", handler.Handle)

	req := httptest.NewRequest("GET", "/api/datasources", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var body struct {
		Datasources []string `json:"datasources"`
		Default     string   `json:"default"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if body.Default != "prod" {
		t.Errorf("expected default 'prod', got %q", body.Default)
	}
	if len(body.Datasources) != 2 {
		t.Fatalf("expected 2 datasources, got %d", len(body.Datasources))
	}
	// Names are sorted
	if body.Datasources[0] != "prod" || body.Datasources[1] != "staging" {
		t.Errorf("expected [prod, staging], got %v", body.Datasources)
	}
}

func TestDatasourcesHandlerSingle(t *testing.T) {
	registry := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "default", Type: "prometheus", URL: "http://localhost:9090", Timeout: 5 * time.Second, Default: true},
	})
	handler := NewDatasourcesHandler(registry)

	router := gin.New()
	router.GET("/api/datasources", handler.Handle)

	req := httptest.NewRequest("GET", "/api/datasources", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var body struct {
		Datasources []string `json:"datasources"`
		Default     string   `json:"default"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if body.Default != "default" {
		t.Errorf("expected default 'default', got %q", body.Default)
	}
	if len(body.Datasources) != 1 {
		t.Errorf("expected 1 datasource, got %d", len(body.Datasources))
	}
}
