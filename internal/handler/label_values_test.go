package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/config"
	"github.com/tokuhirom/dashyard/internal/datasource"
)

func TestLabelValuesHandler(t *testing.T) {
	promServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success","data":["eth0","eth1"]}`))
	}))
	defer promServer.Close()

	registry := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "default", Type: "prometheus", URL: promServer.URL, Timeout: 5 * time.Second, Default: true},
	})
	handler := NewLabelValuesHandler(registry)

	router := gin.New()
	router.GET("/api/label-values", handler.Handle)

	req := httptest.NewRequest("GET", "/api/label-values?label=device&match=system_network_io_bytes_total", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	expected := `{"status":"success","data":["eth0","eth1"]}`
	if resp.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, resp.Body.String())
	}
}

func TestLabelValuesHandlerMissingLabel(t *testing.T) {
	registry := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "default", Type: "prometheus", URL: "http://localhost:9090", Timeout: 5 * time.Second, Default: true},
	})
	handler := NewLabelValuesHandler(registry)

	router := gin.New()
	router.GET("/api/label-values", handler.Handle)

	req := httptest.NewRequest("GET", "/api/label-values", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestLabelValuesHandlerPrometheusError(t *testing.T) {
	promServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"status":"error","error":"internal error"}`))
	}))
	defer promServer.Close()

	registry := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "default", Type: "prometheus", URL: promServer.URL, Timeout: 5 * time.Second, Default: true},
	})
	handler := NewLabelValuesHandler(registry)

	router := gin.New()
	router.GET("/api/label-values", handler.Handle)

	req := httptest.NewRequest("GET", "/api/label-values?label=device", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.Code)
	}
}

func TestLabelValuesHandlerWithDatasourceParam(t *testing.T) {
	promServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":["val1"]}`))
	}))
	defer promServer.Close()

	registry := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "main", Type: "prometheus", URL: "http://localhost:1", Timeout: 5 * time.Second, Default: true},
		{Name: "app", Type: "prometheus", URL: promServer.URL, Timeout: 5 * time.Second},
	})
	handler := NewLabelValuesHandler(registry)

	router := gin.New()
	router.GET("/api/label-values", handler.Handle)

	req := httptest.NewRequest("GET", "/api/label-values?label=device&datasource=app", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}
}
