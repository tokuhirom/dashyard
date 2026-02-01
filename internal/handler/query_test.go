package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/config"
	"github.com/tokuhirom/dashyard/internal/datasource"
)

func TestQueryHandler(t *testing.T) {
	promServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[]}}`))
	}))
	defer promServer.Close()

	registry, _ := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "default", Type: "prometheus", URL: promServer.URL, Timeout: 5 * time.Second, Default: true},
	})
	handler := NewQueryHandler(registry)

	router := gin.New()
	router.GET("/api/query", handler.Handle)

	req := httptest.NewRequest("GET", "/api/query?query=up&start=1000&end=2000&step=15s", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	expected := `{"status":"success","data":{"resultType":"matrix","result":[]}}`
	if resp.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, resp.Body.String())
	}
}

func TestQueryHandlerMissingParams(t *testing.T) {
	registry, _ := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "default", Type: "prometheus", URL: "http://localhost:9090", Timeout: 5 * time.Second, Default: true},
	})
	handler := NewQueryHandler(registry)

	router := gin.New()
	router.GET("/api/query", handler.Handle)

	tests := []struct {
		name string
		url  string
	}{
		{"missing query", "/api/query?start=1000&end=2000&step=15s"},
		{"missing start", "/api/query?query=up&end=2000&step=15s"},
		{"missing end", "/api/query?query=up&start=1000&step=15s"},
		{"missing step", "/api/query?query=up&start=1000&end=2000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", resp.Code)
			}
		})
	}
}

func TestQueryHandlerPrometheusError(t *testing.T) {
	promServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"status":"error","error":"internal error"}`))
	}))
	defer promServer.Close()

	registry, _ := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "default", Type: "prometheus", URL: promServer.URL, Timeout: 5 * time.Second, Default: true},
	})
	handler := NewQueryHandler(registry)

	router := gin.New()
	router.GET("/api/query", handler.Handle)

	req := httptest.NewRequest("GET", "/api/query?query=up&start=1000&end=2000&step=15s", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.Code)
	}
}

func TestQueryHandlerWithDatasourceParam(t *testing.T) {
	mainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"src":"main"}}]}}`))
	}))
	defer mainServer.Close()

	appServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"src":"app"}}]}}`))
	}))
	defer appServer.Close()

	registry, _ := datasource.NewRegistry([]config.DatasourceConfig{
		{Name: "main", Type: "prometheus", URL: mainServer.URL, Timeout: 5 * time.Second, Default: true},
		{Name: "app", Type: "prometheus", URL: appServer.URL, Timeout: 5 * time.Second},
	})
	handler := NewQueryHandler(registry)

	router := gin.New()
	router.GET("/api/query", handler.Handle)

	// Query with explicit datasource=app
	req := httptest.NewRequest("GET", "/api/query?query=up&start=1000&end=2000&step=15s&datasource=app", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}
	if body := resp.Body.String(); !strings.Contains(body, "app") {
		t.Errorf("expected response from app server, got %q", body)
	}

	// Query with unknown datasource
	req = httptest.NewRequest("GET", "/api/query?query=up&start=1000&end=2000&step=15s&datasource=nonexistent", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for unknown datasource, got %d", resp.Code)
	}
}
