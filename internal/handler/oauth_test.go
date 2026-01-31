package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/auth"
	"github.com/tokuhirom/dashyard/internal/config"
)

func TestOAuthLogout(t *testing.T) {
	sm := auth.NewSessionManager("test-secret-that-is-32bytes!!", false)
	handler := NewOAuthHandler(nil, sm)

	router := gin.New()
	router.GET("/auth/logout", handler.Logout)

	req := httptest.NewRequest("GET", "/auth/logout", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusTemporaryRedirect {
		t.Errorf("expected 307, got %d", resp.Code)
	}
	if loc := resp.Header().Get("Location"); loc != "/" {
		t.Errorf("expected redirect to '/', got %q", loc)
	}
}

func TestOAuthBeginAuthUnknownProvider(t *testing.T) {
	sm := auth.NewSessionManager("test-secret-that-is-32bytes!!", false)
	providers := []config.OAuthProviderConfig{}
	handler := NewOAuthHandler(providers, sm)

	router := gin.New()
	router.GET("/auth/:provider", handler.BeginAuth)

	req := httptest.NewRequest("GET", "/auth/unknown", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// gothic will fail with an error for an unknown provider.
	// The response should not be 200.
	if resp.Code == http.StatusOK {
		t.Error("expected non-200 status for unknown provider")
	}
}
