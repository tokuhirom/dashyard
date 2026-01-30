package prometheus

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestQueryRange(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/query_range" {
			t.Errorf("expected path '/api/v1/query_range', got %q", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "up" {
			t.Errorf("expected query 'up', got %q", r.URL.Query().Get("query"))
		}
		if r.URL.Query().Get("start") != "1000" {
			t.Errorf("expected start '1000', got %q", r.URL.Query().Get("start"))
		}
		if r.URL.Query().Get("end") != "2000" {
			t.Errorf("expected end '2000', got %q", r.URL.Query().Get("end"))
		}
		if r.URL.Query().Get("step") != "15s" {
			t.Errorf("expected step '15s', got %q", r.URL.Query().Get("step"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[]}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)

	body, statusCode, err := client.QueryRange(context.Background(), "up", "1000", "2000", "15s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer body.Close()

	if statusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", statusCode)
	}

	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("unexpected error reading body: %v", err)
	}
	expected := `{"status":"success","data":{"resultType":"matrix","result":[]}}`
	if string(data) != expected {
		t.Errorf("expected body %q, got %q", expected, string(data))
	}
}

func TestQueryRangeServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","error":"internal error"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)

	body, statusCode, err := client.QueryRange(context.Background(), "up", "1000", "2000", "15s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer body.Close()

	if statusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", statusCode)
	}
}

func TestQueryRangeConnectionError(t *testing.T) {
	client := NewClient("http://localhost:1", 1*time.Second)

	_, _, err := client.QueryRange(context.Background(), "up", "1000", "2000", "15s")
	if err == nil {
		t.Error("expected error for connection failure")
	}
}
