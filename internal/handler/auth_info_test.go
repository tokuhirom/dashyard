package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/config"
)

func TestAuthInfoPasswordOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	users := []config.User{{ID: "admin", PasswordHash: "hash"}}
	h := NewAuthInfoHandler(users, nil)
	r.GET("/api/auth-info", h.Handle)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/auth-info", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp AuthInfoResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.PasswordEnabled {
		t.Error("expected password_enabled=true")
	}
	if resp.OAuthEnabled {
		t.Error("expected oauth_enabled=false")
	}
	if resp.OAuthProvider != "" {
		t.Errorf("expected empty oauth_provider, got %q", resp.OAuthProvider)
	}
}

func TestAuthInfoOAuthOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	oauthCfg := &config.OAuthConfig{
		Provider: "github",
	}
	h := NewAuthInfoHandler(nil, oauthCfg)
	r.GET("/api/auth-info", h.Handle)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/auth-info", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp AuthInfoResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.PasswordEnabled {
		t.Error("expected password_enabled=false")
	}
	if !resp.OAuthEnabled {
		t.Error("expected oauth_enabled=true")
	}
	if resp.OAuthProvider != "github" {
		t.Errorf("expected oauth_provider 'github', got %q", resp.OAuthProvider)
	}
	if resp.OAuthLoginURL != "/auth/login" {
		t.Errorf("expected oauth_login_url '/auth/login', got %q", resp.OAuthLoginURL)
	}
}

func TestAuthInfoBoth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	users := []config.User{{ID: "admin", PasswordHash: "hash"}}
	oauthCfg := &config.OAuthConfig{
		Provider: "google",
	}
	h := NewAuthInfoHandler(users, oauthCfg)
	r.GET("/api/auth-info", h.Handle)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/auth-info", nil)
	r.ServeHTTP(w, req)

	var resp AuthInfoResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.PasswordEnabled {
		t.Error("expected password_enabled=true")
	}
	if !resp.OAuthEnabled {
		t.Error("expected oauth_enabled=true")
	}
}
