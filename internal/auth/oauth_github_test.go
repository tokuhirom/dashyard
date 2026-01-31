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

func TestGitHubProviderUserInfo(t *testing.T) {
	// Mock GitHub API
	mux := http.NewServeMux()
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"login": "testuser",
			"id":    12345,
			"email": "test@example.com",
		})
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cfg := &config.OAuthConfig{
		Provider:     "github",
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost/callback",
	}

	p := NewGitHubProvider(cfg)
	p.userURL = server.URL + "/user"

	token := &oauth2.Token{AccessToken: "test-token"}
	token.TokenType = "bearer"

	info, err := p.UserInfo(context.Background(), token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %q", info.Username)
	}
	if info.ID != "12345" {
		t.Errorf("expected id '12345', got %q", info.ID)
	}
	if info.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", info.Email)
	}
}

func TestGitHubProviderUserInfoWithOrgs(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"login": "testuser",
			"id":    12345,
			"email": "test@example.com",
		})
	})
	mux.HandleFunc("/user/orgs", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{
			{"login": "org-one"},
			{"login": "org-two"},
		})
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cfg := &config.OAuthConfig{
		Provider:     "github",
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost/callback",
		AllowedOrgs:  []string{"org-one"},
	}

	p := NewGitHubProvider(cfg)
	p.userURL = server.URL + "/user"
	p.orgsURL = server.URL + "/user/orgs"

	token := &oauth2.Token{AccessToken: "test-token"}
	token.TokenType = "bearer"

	info, err := p.UserInfo(context.Background(), token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(info.Orgs) != 2 {
		t.Fatalf("expected 2 orgs, got %d", len(info.Orgs))
	}
	if info.Orgs[0] != "org-one" {
		t.Errorf("expected first org 'org-one', got %q", info.Orgs[0])
	}
}

func TestGitHubProviderDefaultScopes(t *testing.T) {
	cfg := &config.OAuthConfig{
		Provider:     "github",
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost/callback",
	}

	p := NewGitHubProvider(cfg)
	if len(p.config.Scopes) != 2 {
		t.Fatalf("expected 2 default scopes, got %d", len(p.config.Scopes))
	}
	if p.config.Scopes[0] != "read:user" || p.config.Scopes[1] != "read:org" {
		t.Errorf("unexpected default scopes: %v", p.config.Scopes)
	}
}

func TestGitHubProviderCustomScopes(t *testing.T) {
	cfg := &config.OAuthConfig{
		Provider:     "github",
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost/callback",
		Scopes:       []string{"user"},
	}

	p := NewGitHubProvider(cfg)
	if len(p.config.Scopes) != 1 || p.config.Scopes[0] != "user" {
		t.Errorf("expected scopes [user], got %v", p.config.Scopes)
	}
}

func TestGitHubProviderName(t *testing.T) {
	p := &GitHubProvider{}
	if p.Name() != "GitHub" {
		t.Errorf("expected name 'GitHub', got %q", p.Name())
	}
}

func TestGitHubProviderAuthCodeURL(t *testing.T) {
	cfg := &config.OAuthConfig{
		Provider:     "github",
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost/callback",
	}

	p := NewGitHubProvider(cfg)
	url := p.AuthCodeURL("test-state")
	if url == "" {
		t.Fatal("expected non-empty auth URL")
	}
}
