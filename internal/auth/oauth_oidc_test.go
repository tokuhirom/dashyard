package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tokuhirom/dashyard/internal/config"
	"golang.org/x/oauth2"
)

func setupOIDCServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		// We need server URL, but we don't have it yet. Use relative paths resolved by the client.
		// Actually, the discovery doc should have absolute URLs. We'll fix this in the test.
	})
	return httptest.NewServer(mux)
}

func TestOIDCProviderDiscoveryAndUserInfo(t *testing.T) {
	var serverURL string
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"authorization_endpoint": serverURL + "/authorize",
			"token_endpoint":         serverURL + "/token",
			"userinfo_endpoint":      serverURL + "/userinfo",
		})
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"sub":                "user-123",
			"preferred_username": "testuser",
			"email":              "test@example.com",
		})
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	serverURL = server.URL

	cfg := &config.OAuthConfig{
		Provider:     "oidc",
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost/callback",
		IssuerURL:    server.URL,
	}

	p := NewOIDCProvider(cfg, server.URL)
	p.httpClient = server.Client()

	// Test discovery by getting auth URL
	url := p.AuthCodeURL("test-state")
	if url == "" {
		t.Fatal("expected non-empty auth URL")
	}

	// Test userinfo
	token := &oauth2.Token{AccessToken: "test-token"}
	token.TokenType = "bearer"

	info, err := p.UserInfo(context.Background(), token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ID != "user-123" {
		t.Errorf("expected ID 'user-123', got %q", info.ID)
	}
	if info.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %q", info.Username)
	}
	if info.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", info.Email)
	}
}

func TestOIDCProviderFallbackUsername(t *testing.T) {
	var serverURL string
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"authorization_endpoint": serverURL + "/authorize",
			"token_endpoint":         serverURL + "/token",
			"userinfo_endpoint":      serverURL + "/userinfo",
		})
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"sub":   "user-456",
			"email": "user@example.com",
		})
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	serverURL = server.URL

	cfg := &config.OAuthConfig{
		Provider:     "oidc",
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost/callback",
		IssuerURL:    server.URL,
	}

	p := NewOIDCProvider(cfg, server.URL)
	p.httpClient = server.Client()

	token := &oauth2.Token{AccessToken: "test-token"}
	token.TokenType = "bearer"

	info, err := p.UserInfo(context.Background(), token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should fall back to email when preferred_username is empty
	if info.Username != "user@example.com" {
		t.Errorf("expected username 'user@example.com', got %q", info.Username)
	}
}

func TestOIDCProviderName(t *testing.T) {
	cfg := &config.OAuthConfig{Provider: "oidc"}
	p := NewOIDCProvider(cfg, "https://example.com")
	if p.Name() != "OIDC" {
		t.Errorf("expected name 'OIDC', got %q", p.Name())
	}
}

func TestOIDCProviderGoogleName(t *testing.T) {
	cfg := &config.OAuthConfig{Provider: "google"}
	p := NewOIDCProvider(cfg, "https://accounts.google.com")
	if p.Name() != "Google" {
		t.Errorf("expected name 'Google', got %q", p.Name())
	}
}

func TestOIDCProviderDefaultScopes(t *testing.T) {
	var serverURL string
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"authorization_endpoint": serverURL + "/authorize",
			"token_endpoint":         serverURL + "/token",
			"userinfo_endpoint":      serverURL + "/userinfo",
		})
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	serverURL = server.URL

	cfg := &config.OAuthConfig{
		Provider:     "oidc",
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost/callback",
		IssuerURL:    server.URL,
	}

	p := NewOIDCProvider(cfg, server.URL)
	p.httpClient = server.Client()

	oc, err := p.oauthConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(oc.Scopes) != 3 {
		t.Fatalf("expected 3 default scopes, got %d", len(oc.Scopes))
	}
}
