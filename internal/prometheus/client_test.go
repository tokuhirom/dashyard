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
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[]}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)

	body, statusCode, err := client.QueryRange(context.Background(), "up", "1000", "2000", "15s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = body.Close() }()

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
		_, _ = w.Write([]byte(`{"status":"error","error":"internal error"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)

	body, statusCode, err := client.QueryRange(context.Background(), "up", "1000", "2000", "15s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = body.Close() }()

	if statusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", statusCode)
	}
}

func TestLabelValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/label/device/values" {
			t.Errorf("expected path '/api/v1/label/device/values', got %q", r.URL.Path)
		}
		if r.URL.Query().Get("match[]") != "system_network_io_bytes_total" {
			t.Errorf("expected match[] 'system_network_io_bytes_total', got %q", r.URL.Query().Get("match[]"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success","data":["eth0","eth1"]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)

	body, statusCode, err := client.LabelValues(context.Background(), "device", "system_network_io_bytes_total")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = body.Close() }()

	if statusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", statusCode)
	}

	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("unexpected error reading body: %v", err)
	}
	expected := `{"status":"success","data":["eth0","eth1"]}`
	if string(data) != expected {
		t.Errorf("expected body %q, got %q", expected, string(data))
	}
}

func TestLabelValuesNoMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/label/cpu/values" {
			t.Errorf("expected path '/api/v1/label/cpu/values', got %q", r.URL.Path)
		}
		if r.URL.Query().Get("match[]") != "" {
			t.Errorf("expected no match[] parameter, got %q", r.URL.Query().Get("match[]"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success","data":["cpu0","cpu1"]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)

	body, statusCode, err := client.LabelValues(context.Background(), "cpu", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = body.Close() }()

	if statusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", statusCode)
	}
}

func TestPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/-/ready" {
			t.Errorf("expected path '/-/ready', got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)

	if err := client.Ping(context.Background()); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestPingNotReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)

	if err := client.Ping(context.Background()); err == nil {
		t.Error("expected error for non-200 response")
	}
}

func TestPingConnectionError(t *testing.T) {
	client := NewClient("http://localhost:1", 1*time.Second)

	if err := client.Ping(context.Background()); err == nil {
		t.Error("expected error for connection failure")
	}
}

func TestQueryRangeWithBasePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/prod/api/v1/query_range" {
			t.Errorf("expected path '/prod/api/v1/query_range', got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[]}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL+"/prod", 5*time.Second)

	body, statusCode, err := client.QueryRange(context.Background(), "up", "1000", "2000", "15s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = body.Close() }()

	if statusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", statusCode)
	}
}

func TestPingWithBasePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/prod/-/ready" {
			t.Errorf("expected path '/prod/-/ready', got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL+"/prod", 5*time.Second)

	if err := client.Ping(context.Background()); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestLabelValuesWithBasePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/prod/api/v1/label/device/values" {
			t.Errorf("expected path '/prod/api/v1/label/device/values', got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success","data":["eth0"]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL+"/prod", 5*time.Second)

	body, statusCode, err := client.LabelValues(context.Background(), "device", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = body.Close() }()

	if statusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", statusCode)
	}
}

func TestQueryRangeConnectionError(t *testing.T) {
	client := NewClient("http://localhost:1", 1*time.Second)

	_, _, err := client.QueryRange(context.Background(), "up", "1000", "2000", "15s")
	if err == nil {
		t.Error("expected error for connection failure")
	}
}
