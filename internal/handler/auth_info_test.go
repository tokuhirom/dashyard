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
	users := []config.User{{ID: "admin", PasswordHash: "hash"}}
	handler := NewAuthInfoHandler(users, nil)

	router := gin.New()
	router.GET("/api/auth-info", handler.Handle)

	req := httptest.NewRequest("GET", "/api/auth-info", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var result AuthInfoResponse
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if !result.PasswordEnabled {
		t.Error("expected password_enabled=true")
	}
	if len(result.OAuthProviders) != 0 {
		t.Errorf("expected 0 oauth providers, got %d", len(result.OAuthProviders))
	}
}

func TestAuthInfoOAuthOnly(t *testing.T) {
	providers := []config.OAuthProviderConfig{
		{Provider: "github", ClientID: "id", ClientSecret: "secret"},
	}
	handler := NewAuthInfoHandler(nil, providers)

	router := gin.New()
	router.GET("/api/auth-info", handler.Handle)

	req := httptest.NewRequest("GET", "/api/auth-info", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var result AuthInfoResponse
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.PasswordEnabled {
		t.Error("expected password_enabled=false")
	}
	if len(result.OAuthProviders) != 1 {
		t.Fatalf("expected 1 oauth provider, got %d", len(result.OAuthProviders))
	}
	if result.OAuthProviders[0].Name != "github" {
		t.Errorf("expected provider name 'github', got %q", result.OAuthProviders[0].Name)
	}
	if result.OAuthProviders[0].URL != "/auth/github" {
		t.Errorf("expected url '/auth/github', got %q", result.OAuthProviders[0].URL)
	}
}

func TestAuthInfoBothMethods(t *testing.T) {
	users := []config.User{{ID: "admin", PasswordHash: "hash"}}
	providers := []config.OAuthProviderConfig{
		{Provider: "github", ClientID: "id", ClientSecret: "secret"},
	}
	handler := NewAuthInfoHandler(users, providers)

	router := gin.New()
	router.GET("/api/auth-info", handler.Handle)

	req := httptest.NewRequest("GET", "/api/auth-info", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var result AuthInfoResponse
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if !result.PasswordEnabled {
		t.Error("expected password_enabled=true")
	}
	if len(result.OAuthProviders) != 1 {
		t.Errorf("expected 1 oauth provider, got %d", len(result.OAuthProviders))
	}
}
