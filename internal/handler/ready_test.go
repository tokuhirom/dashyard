package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/prometheus"
)

func TestReadyHandler_OK(t *testing.T) {
	promServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/-/ready" {
			t.Errorf("expected path '/-/ready', got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer promServer.Close()

	client := prometheus.NewClient(promServer.URL, 5*time.Second)
	h := NewReadyHandler(client)

	router := gin.New()
	router.GET("/ready", h.Handle)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", body["status"])
	}
	if body["prometheus"] != "reachable" {
		t.Errorf("expected prometheus 'reachable', got %q", body["prometheus"])
	}
}

func TestReadyHandler_PrometheusUnreachable(t *testing.T) {
	client := prometheus.NewClient("http://localhost:1", 1*time.Second)
	h := NewReadyHandler(client)

	router := gin.New()
	router.GET("/ready", h.Handle)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", resp.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if body["status"] != "degraded" {
		t.Errorf("expected status 'degraded', got %q", body["status"])
	}
	if body["prometheus"] != "unreachable" {
		t.Errorf("expected prometheus 'unreachable', got %q", body["prometheus"])
	}
}

func TestReadyHandler_PrometheusNotReady(t *testing.T) {
	promServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer promServer.Close()

	client := prometheus.NewClient(promServer.URL, 5*time.Second)
	h := NewReadyHandler(client)

	router := gin.New()
	router.GET("/ready", h.Handle)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", resp.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if body["status"] != "degraded" {
		t.Errorf("expected status 'degraded', got %q", body["status"])
	}
	if body["prometheus"] != "unreachable" {
		t.Errorf("expected prometheus 'unreachable', got %q", body["prometheus"])
	}
}
