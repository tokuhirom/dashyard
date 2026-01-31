package prometheus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMetricNames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/label/__name__/values" {
			t.Errorf("expected path '/api/v1/label/__name__/values', got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":["go_gc_duration_seconds","up"]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	names, err := client.MetricNames(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
	if names[0] != "go_gc_duration_seconds" {
		t.Errorf("expected first name 'go_gc_duration_seconds', got %q", names[0])
	}
	if names[1] != "up" {
		t.Errorf("expected second name 'up', got %q", names[1])
	}
}

func TestMetricMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/metadata" {
			t.Errorf("expected path '/api/v1/metadata', got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":{"up":[{"type":"gauge","help":"Whether the target is up.","unit":""}]}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	metadata, err := client.MetricMetadata(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, ok := metadata["up"]
	if !ok {
		t.Fatal("expected metadata for 'up'")
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Type != "gauge" {
		t.Errorf("expected type 'gauge', got %q", entries[0].Type)
	}
	if entries[0].Help != "Whether the target is up." {
		t.Errorf("expected help 'Whether the target is up.', got %q", entries[0].Help)
	}
}

func TestMetricLabels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/labels" {
			t.Errorf("expected path '/api/v1/labels', got %q", r.URL.Path)
		}
		match := r.URL.Query().Get("match[]")
		if match != `{__name__="up"}` {
			t.Errorf("expected match[] '{__name__=\"up\"}', got %q", match)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":["__name__","instance","job"]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	labels, err := client.MetricLabels(context.Background(), "up")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels (without __name__), got %d: %v", len(labels), labels)
	}
	if labels[0] != "instance" {
		t.Errorf("expected first label 'instance', got %q", labels[0])
	}
	if labels[1] != "job" {
		t.Errorf("expected second label 'job', got %q", labels[1])
	}
}

func TestBearerTokenSent(t *testing.T) {
	var receivedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":["up"]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second, WithBearerToken("test-token-123"))
	_, err := client.MetricNames(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "Bearer test-token-123"
	if receivedAuth != expected {
		t.Errorf("expected Authorization %q, got %q", expected, receivedAuth)
	}
}

func TestBearerTokenNotSentWhenEmpty(t *testing.T) {
	var receivedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":["up"]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	_, err := client.MetricNames(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedAuth != "" {
		t.Errorf("expected no Authorization header, got %q", receivedAuth)
	}
}

func TestDoGetHandlesPrometheusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"status":"error","error":"bad_data: invalid label name"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	_, err := client.MetricNames(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); got == "" {
		t.Error("expected non-empty error message")
	}
}
