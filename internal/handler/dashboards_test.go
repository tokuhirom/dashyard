package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/dashboard"
)

func loadTestHolder(t *testing.T) *dashboard.StoreHolder {
	t.Helper()
	store, err := dashboard.LoadDir("../../testdata/dashboards")
	if err != nil {
		t.Fatalf("failed to load test dashboards: %v", err)
	}
	return dashboard.NewStoreHolder(store)
}

func TestDashboardsList(t *testing.T) {
	holder := loadTestHolder(t)
	handler := NewDashboardsHandler(holder, "Dashyard", "")

	router := gin.New()
	router.GET("/api/dashboards", handler.List)

	req := httptest.NewRequest("GET", "/api/dashboards", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var result struct {
		Dashboards []struct {
			Path  string `json:"path"`
			Title string `json:"title"`
		} `json:"dashboards"`
		Tree json.RawMessage `json:"tree"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(result.Dashboards) != 3 {
		t.Errorf("expected 3 dashboards, got %d", len(result.Dashboards))
	}
}

func TestDashboardsGetDeepNested(t *testing.T) {
	holder := loadTestHolder(t)
	handler := NewDashboardsHandler(holder, "Dashyard", "")

	router := gin.New()
	router.GET("/api/dashboards/*path", handler.Get)

	req := httptest.NewRequest("GET", "/api/dashboards/infra/cloud/sakura", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result struct {
		Title string `json:"title"`
		Path  string `json:"path"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Title != "Sakura Metrics" {
		t.Errorf("expected title 'Sakura Metrics', got %q", result.Title)
	}
	if result.Path != "infra/cloud/sakura" {
		t.Errorf("expected path 'infra/cloud/sakura', got %q", result.Path)
	}
}

func TestDashboardsGet(t *testing.T) {
	holder := loadTestHolder(t)
	handler := NewDashboardsHandler(holder, "Dashyard", "")

	router := gin.New()
	router.GET("/api/dashboards/*path", handler.Get)

	req := httptest.NewRequest("GET", "/api/dashboards/overview", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result struct {
		Title string `json:"title"`
		Path  string `json:"path"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Title != "System Overview" {
		t.Errorf("expected title 'System Overview', got %q", result.Title)
	}
	if result.Path != "overview" {
		t.Errorf("expected path 'overview', got %q", result.Path)
	}
}

func TestDashboardsGetNested(t *testing.T) {
	holder := loadTestHolder(t)
	handler := NewDashboardsHandler(holder, "Dashyard", "")

	router := gin.New()
	router.GET("/api/dashboards/*path", handler.Get)

	req := httptest.NewRequest("GET", "/api/dashboards/infra/network", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result struct {
		Title string `json:"title"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Title != "Network" {
		t.Errorf("expected title 'Network', got %q", result.Title)
	}
}

func TestDashboardsGetNotFound(t *testing.T) {
	holder := loadTestHolder(t)
	handler := NewDashboardsHandler(holder, "Dashyard", "")

	router := gin.New()
	router.GET("/api/dashboards/*path", handler.Get)

	req := httptest.NewRequest("GET", "/api/dashboards/nonexistent", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.Code)
	}
}
