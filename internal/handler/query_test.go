package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/prometheus"
)

func TestQueryHandler(t *testing.T) {
	promServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[]}}`))
	}))
	defer promServer.Close()

	client := prometheus.NewClient(promServer.URL, 5*time.Second)
	handler := NewQueryHandler(client)

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
	client := prometheus.NewClient("http://localhost:9090", 5*time.Second)
	handler := NewQueryHandler(client)

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

	client := prometheus.NewClient(promServer.URL, 5*time.Second)
	handler := NewQueryHandler(client)

	router := gin.New()
	router.GET("/api/query", handler.Handle)

	req := httptest.NewRequest("GET", "/api/query?query=up&start=1000&end=2000&step=15s", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.Code)
	}
}
