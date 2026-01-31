package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/auth"
	"github.com/tokuhirom/dashyard/internal/config"
	"golang.org/x/oauth2"
)

// mockOAuthProvider implements auth.OAuthProvider for testing.
type mockOAuthProvider struct {
	authURL  string
	token    *oauth2.Token
	tokenErr error
	info     *auth.OAuthUserInfo
	infoErr  error
}

func (m *mockOAuthProvider) AuthCodeURL(state string) string {
	return m.authURL + "?state=" + state
}

func (m *mockOAuthProvider) Exchange(_ context.Context, _ string) (*oauth2.Token, error) {
	return m.token, m.tokenErr
}

func (m *mockOAuthProvider) UserInfo(_ context.Context, _ *oauth2.Token) (*auth.OAuthUserInfo, error) {
	return m.info, m.infoErr
}

func (m *mockOAuthProvider) Name() string {
	return "Mock"
}

func TestOAuthLoginRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sm := auth.NewSessionManager("test-secret", false)
	stateManager := auth.NewOAuthStateManager("test-secret", false)
	provider := &mockOAuthProvider{authURL: "https://provider.example.com/auth"}
	oauthCfg := &config.OAuthConfig{Provider: "github"}

	h := NewOAuthHandler(provider, stateManager, sm, oauthCfg)

	r := gin.New()
	r.GET("/auth/login", h.Login)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/login", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", w.Code)
	}

	loc := w.Header().Get("Location")
	if loc == "" {
		t.Fatal("expected Location header")
	}
}

func TestOAuthCallbackSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sm := auth.NewSessionManager("test-secret", false)
	stateManager := auth.NewOAuthStateManager("test-secret", false)
	provider := &mockOAuthProvider{
		authURL: "https://provider.example.com/auth",
		token:   &oauth2.Token{AccessToken: "test-token"},
		info: &auth.OAuthUserInfo{
			ID:       "123",
			Username: "testuser",
			Email:    "test@example.com",
		},
	}
	oauthCfg := &config.OAuthConfig{Provider: "github"}
	h := NewOAuthHandler(provider, stateManager, sm, oauthCfg)

	// Generate state first
	stateW := httptest.NewRecorder()
	state, err := stateManager.Generate(stateW)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stateCookie := stateW.Result().Cookies()[0]

	r := gin.New()
	r.GET("/auth/callback", h.Callback)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/callback?code=test-code&state="+state, nil)
	req.AddCookie(stateCookie)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/" {
		t.Errorf("expected redirect to '/', got %q", loc)
	}

	// Verify session cookie was set
	var sessionCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "dashyard_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Error("expected session cookie to be set")
	}
}

func TestOAuthCallbackInvalidState(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sm := auth.NewSessionManager("test-secret", false)
	stateManager := auth.NewOAuthStateManager("test-secret", false)
	provider := &mockOAuthProvider{}
	oauthCfg := &config.OAuthConfig{Provider: "github"}
	h := NewOAuthHandler(provider, stateManager, sm, oauthCfg)

	r := gin.New()
	r.GET("/auth/callback", h.Callback)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/callback?code=test-code&state=invalid", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestOAuthCallbackUserNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sm := auth.NewSessionManager("test-secret", false)
	stateManager := auth.NewOAuthStateManager("test-secret", false)
	provider := &mockOAuthProvider{
		token: &oauth2.Token{AccessToken: "test-token"},
		info: &auth.OAuthUserInfo{
			ID:       "123",
			Username: "evil-user",
		},
	}
	oauthCfg := &config.OAuthConfig{
		Provider:     "github",
		AllowedUsers: []string{"good-user"},
	}
	h := NewOAuthHandler(provider, stateManager, sm, oauthCfg)

	// Generate state
	stateW := httptest.NewRecorder()
	state, err := stateManager.Generate(stateW)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stateCookie := stateW.Result().Cookies()[0]

	r := gin.New()
	r.GET("/auth/callback", h.Callback)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/callback?code=test-code&state="+state, nil)
	req.AddCookie(stateCookie)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestOAuthLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sm := auth.NewSessionManager("test-secret", false)
	stateManager := auth.NewOAuthStateManager("test-secret", false)
	provider := &mockOAuthProvider{}
	oauthCfg := &config.OAuthConfig{Provider: "github"}
	h := NewOAuthHandler(provider, stateManager, sm, oauthCfg)

	r := gin.New()
	r.GET("/auth/logout", h.Logout)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/logout", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", w.Code)
	}

	// Verify session cookie was cleared
	var sessionCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "dashyard_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil || sessionCookie.MaxAge != -1 {
		t.Error("expected session cookie to be cleared")
	}
}
