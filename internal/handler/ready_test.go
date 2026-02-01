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

func TestReadyHandler_OK(t *testing.T) {
	promServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/-/ready" {
			t.Errorf("expected path '/-/ready', got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer promServer.Close()

	registry := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "default", Type: "prometheus", URL: promServer.URL, Timeout: 5 * time.Second, Default: true},
	})
	h := NewReadyHandler(registry)

	router := gin.New()
	router.GET("/ready", h.Handle)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", body["status"])
	}
	ds, ok := body["datasources"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected datasources map, got %T", body["datasources"])
	}
	if ds["default"] != "reachable" {
		t.Errorf("expected datasource 'default' to be 'reachable', got %q", ds["default"])
	}
}

func TestReadyHandler_PrometheusUnreachable(t *testing.T) {
	registry := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "default", Type: "prometheus", URL: "http://localhost:1", Timeout: 1 * time.Second, Default: true},
	})
	h := NewReadyHandler(registry)

	router := gin.New()
	router.GET("/ready", h.Handle)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", resp.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if body["status"] != "degraded" {
		t.Errorf("expected status 'degraded', got %q", body["status"])
	}
	ds, ok := body["datasources"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected datasources map, got %T", body["datasources"])
	}
	if ds["default"] != "unreachable" {
		t.Errorf("expected datasource 'default' to be 'unreachable', got %q", ds["default"])
	}
}

func TestReadyHandler_PrometheusNotReady(t *testing.T) {
	promServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer promServer.Close()

	registry := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "default", Type: "prometheus", URL: promServer.URL, Timeout: 5 * time.Second, Default: true},
	})
	h := NewReadyHandler(registry)

	router := gin.New()
	router.GET("/ready", h.Handle)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", resp.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if body["status"] != "degraded" {
		t.Errorf("expected status 'degraded', got %q", body["status"])
	}
	ds, ok := body["datasources"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected datasources map, got %T", body["datasources"])
	}
	if ds["default"] != "unreachable" {
		t.Errorf("expected datasource 'default' to be 'unreachable', got %q", ds["default"])
	}
}

func TestReadyHandler_MultipleDatasources(t *testing.T) {
	goodServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer goodServer.Close()

	registry := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "good", Type: "prometheus", URL: goodServer.URL, Timeout: 5 * time.Second, Default: true},
		{Name: "bad", Type: "prometheus", URL: "http://localhost:1", Timeout: 1 * time.Second},
	})
	h := NewReadyHandler(registry)

	router := gin.New()
	router.GET("/ready", h.Handle)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Should be degraded because one datasource is unreachable
	if resp.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", resp.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	ds, ok := body["datasources"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected datasources map, got %T", body["datasources"])
	}
	if ds["good"] != "reachable" {
		t.Errorf("expected 'good' to be reachable, got %q", ds["good"])
	}
	if ds["bad"] != "unreachable" {
		t.Errorf("expected 'bad' to be unreachable, got %q", ds["bad"])
	}
}
