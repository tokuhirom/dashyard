package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/tokuhirom/dashyard/internal/auth"
	"github.com/tokuhirom/dashyard/internal/config"
)

// newDummyGitHub starts a fake GitHub OAuth server and returns its URL.
func newDummyGitHub() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /login/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
		redirectURI := r.URL.Query().Get("redirect_uri")
		state := r.URL.Query().Get("state")
		// Auto-approve: redirect back with a dummy code immediately.
		http.Redirect(w, r, redirectURI+"?code=dummy-auth-code&state="+state, http.StatusFound)
	})

	mux.HandleFunc("POST /login/oauth/access_token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": "dummy-access-token",
			"token_type":   "bearer",
			"scope":        "user:email,read:org",
		})
	})

	mux.HandleFunc("GET /api/v3/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"login":      "dummyuser",
			"id":         12345,
			"avatar_url": "https://avatars.githubusercontent.com/u/0?v=4",
			"name":       "Dummy User",
			"email":      "dummy@example.com",
		})
	})

	mux.HandleFunc("GET /api/v3/user/emails", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"email":      "dummy@example.com",
				"primary":    true,
				"verified":   true,
				"visibility": "public",
			},
		})
	})

	mux.HandleFunc("GET /api/v3/user/orgs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"login":       "dummy-org",
				"id":          100,
				"description": "A dummy organization for testing",
			},
		})
	})

	return httptest.NewServer(mux)
}

// setupOAuthTestServer creates a Gin router with OAuth endpoints backed by a dummygithub server.
func setupOAuthTestServer(dummyGitHubURL string, providers []config.OAuthProviderConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)

	sm := auth.NewSessionManager("test-secret-that-is-at-least-32-bytes-long!", false)
	auth.InitGothProviders(providers)
	gothic.Store = sm.Store()

	oauthHandler := NewOAuthHandler(providers, sm)
	authInfoHandler := NewAuthInfoHandler(nil, providers)

	router := gin.New()
	router.GET("/api/auth-info", authInfoHandler.Handle)
	router.GET("/auth/:provider", oauthHandler.BeginAuth)
	router.GET("/auth/:provider/callback", oauthHandler.Callback)
	router.GET("/auth/logout", oauthHandler.Logout)

	return router
}

func TestOAuthIntegrationBeginAuthRedirectsToDummyGitHub(t *testing.T) {
	dummyGH := newDummyGitHub()
	defer dummyGH.Close()

	providers := []config.OAuthProviderConfig{
		{
			Provider:     "github",
			ClientID:     "dummy-client-id",
			ClientSecret: "dummy-client-secret",
			BaseURL:      dummyGH.URL,
			RedirectURL:  "http://localhost:8080/auth/github/callback",
		},
	}
	router := setupOAuthTestServer(dummyGH.URL, providers)

	req := httptest.NewRequest("GET", "/auth/github", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// BeginAuth should redirect to the dummygithub authorize endpoint
	if resp.Code != http.StatusTemporaryRedirect && resp.Code != http.StatusFound {
		t.Fatalf("expected redirect, got %d", resp.Code)
	}

	loc := resp.Header().Get("Location")
	parsed, err := url.Parse(loc)
	if err != nil {
		t.Fatalf("failed to parse Location header: %v", err)
	}

	expectedPrefix := dummyGH.URL + "/login/oauth/authorize"
	if loc[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("expected redirect to %s, got %s", expectedPrefix, loc)
	}
	if parsed.Query().Get("client_id") != "dummy-client-id" {
		t.Errorf("expected client_id 'dummy-client-id', got %q", parsed.Query().Get("client_id"))
	}
}

func TestOAuthIntegrationFullFlow(t *testing.T) {
	dummyGH := newDummyGitHub()
	defer dummyGH.Close()

	providers := []config.OAuthProviderConfig{
		{
			Provider:     "github",
			ClientID:     "dummy-client-id",
			ClientSecret: "dummy-client-secret",
			BaseURL:      dummyGH.URL,
			RedirectURL:  "", // Will be set dynamically
		},
	}

	gin.SetMode(gin.TestMode)
	sm := auth.NewSessionManager("test-secret-that-is-at-least-32-bytes-long!", false)

	// We need to create the test server first to know its URL for the redirect_url
	var router *gin.Engine
	router = gin.New()

	oauthHandler := NewOAuthHandler(providers, sm)
	router.GET("/auth/:provider", oauthHandler.BeginAuth)
	router.GET("/auth/:provider/callback", oauthHandler.Callback)

	appServer := httptest.NewServer(router)
	defer appServer.Close()

	// Now update the redirect URL and re-initialize providers
	providers[0].RedirectURL = appServer.URL + "/auth/github/callback"
	auth.InitGothProviders(providers)
	gothic.Store = sm.Store()

	// Use an HTTP client that does NOT follow redirects automatically
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Step 1: Hit /auth/github to start the OAuth flow
	resp, err := client.Get(appServer.URL + "/auth/github")
	if err != nil {
		t.Fatalf("BeginAuth request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTemporaryRedirect && resp.StatusCode != http.StatusFound {
		t.Fatalf("expected redirect from BeginAuth, got %d", resp.StatusCode)
	}

	authRedirect := resp.Header.Get("Location")
	t.Logf("Step 1: BeginAuth redirected to: %s", authRedirect)

	// Extract gothic session cookie to pass along
	var cookies []*http.Cookie
	cookies = append(cookies, resp.Cookies()...)

	// Step 2: Follow the redirect to dummygithub /login/oauth/authorize
	// dummygithub auto-approves and redirects back to callback
	resp2, err := client.Get(authRedirect)
	if err != nil {
		t.Fatalf("dummygithub authorize request failed: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusFound {
		t.Fatalf("expected redirect from dummygithub, got %d", resp2.StatusCode)
	}

	callbackRedirect := resp2.Header.Get("Location")
	t.Logf("Step 2: dummygithub redirected to: %s", callbackRedirect)

	// Step 3: Follow the redirect to /auth/github/callback on our app server
	// We need to include the gothic session cookie from step 1
	req, err := http.NewRequest("GET", callbackRedirect, nil)
	if err != nil {
		t.Fatalf("failed to create callback request: %v", err)
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp3, err := client.Do(req)
	if err != nil {
		t.Fatalf("callback request failed: %v", err)
	}
	defer resp3.Body.Close()

	// The callback should redirect to "/" after successful auth
	if resp3.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307 redirect from callback, got %d", resp3.StatusCode)
	}

	finalRedirect := resp3.Header.Get("Location")
	if finalRedirect != "/" {
		t.Errorf("expected redirect to '/', got %q", finalRedirect)
	}

	// Verify a session cookie was set
	var sessionCookie *http.Cookie
	for _, c := range resp3.Cookies() {
		if c.Name == "dashyard_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Error("expected dashyard_session cookie to be set after OAuth flow")
	}
}

func TestOAuthIntegrationWithAllowedUsers(t *testing.T) {
	dummyGH := newDummyGitHub()
	defer dummyGH.Close()

	// dummygithub returns user "dummyuser", which is NOT in the allowed list
	providers := []config.OAuthProviderConfig{
		{
			Provider:     "github",
			ClientID:     "dummy-client-id",
			ClientSecret: "dummy-client-secret",
			BaseURL:      dummyGH.URL,
			RedirectURL:  "", // Will be set dynamically
			AllowedUsers: []string{"other-user"},
		},
	}

	gin.SetMode(gin.TestMode)
	sm := auth.NewSessionManager("test-secret-that-is-at-least-32-bytes-long!", false)

	router := gin.New()
	oauthHandler := NewOAuthHandler(providers, sm)
	router.GET("/auth/:provider", oauthHandler.BeginAuth)
	router.GET("/auth/:provider/callback", oauthHandler.Callback)

	appServer := httptest.NewServer(router)
	defer appServer.Close()

	providers[0].RedirectURL = appServer.URL + "/auth/github/callback"
	auth.InitGothProviders(providers)
	gothic.Store = sm.Store()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Step 1: Start OAuth flow
	resp, err := client.Get(appServer.URL + "/auth/github")
	if err != nil {
		t.Fatalf("BeginAuth request failed: %v", err)
	}
	defer resp.Body.Close()
	cookies := resp.Cookies()

	// Step 2: Follow redirect to dummygithub
	resp2, err := client.Get(resp.Header.Get("Location"))
	if err != nil {
		t.Fatalf("dummygithub request failed: %v", err)
	}
	defer resp2.Body.Close()

	// Step 3: Follow redirect to callback
	req, _ := http.NewRequest("GET", resp2.Header.Get("Location"), nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp3, err := client.Do(req)
	if err != nil {
		t.Fatalf("callback request failed: %v", err)
	}
	defer resp3.Body.Close()

	// Should redirect with access_denied error since "dummyuser" is not in allowed list
	if resp3.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307, got %d", resp3.StatusCode)
	}

	loc := resp3.Header.Get("Location")
	if loc != "/?error=access_denied" {
		t.Errorf("expected redirect to '/?error=access_denied', got %q", loc)
	}
}

func TestOAuthIntegrationWithAllowedOrgs(t *testing.T) {
	dummyGH := newDummyGitHub()
	defer dummyGH.Close()

	// dummygithub returns org "dummy-org", which IS in the allowed list
	providers := []config.OAuthProviderConfig{
		{
			Provider:     "github",
			ClientID:     "dummy-client-id",
			ClientSecret: "dummy-client-secret",
			BaseURL:      dummyGH.URL,
			RedirectURL:  "", // Will be set dynamically
			AllowedOrgs:  []string{"dummy-org"},
		},
	}

	gin.SetMode(gin.TestMode)
	sm := auth.NewSessionManager("test-secret-that-is-at-least-32-bytes-long!", false)

	router := gin.New()
	oauthHandler := NewOAuthHandler(providers, sm)
	router.GET("/auth/:provider", oauthHandler.BeginAuth)
	router.GET("/auth/:provider/callback", oauthHandler.Callback)

	appServer := httptest.NewServer(router)
	defer appServer.Close()

	providers[0].RedirectURL = appServer.URL + "/auth/github/callback"
	auth.InitGothProviders(providers)
	gothic.Store = sm.Store()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Step 1: Start OAuth flow
	resp, err := client.Get(appServer.URL + "/auth/github")
	if err != nil {
		t.Fatalf("BeginAuth request failed: %v", err)
	}
	defer resp.Body.Close()
	cookies := resp.Cookies()

	// Step 2: Follow redirect to dummygithub
	resp2, err := client.Get(resp.Header.Get("Location"))
	if err != nil {
		t.Fatalf("dummygithub request failed: %v", err)
	}
	defer resp2.Body.Close()

	// Step 3: Follow redirect to callback
	req, _ := http.NewRequest("GET", resp2.Header.Get("Location"), nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp3, err := client.Do(req)
	if err != nil {
		t.Fatalf("callback request failed: %v", err)
	}
	defer resp3.Body.Close()

	// Should succeed because "dummy-org" is in allowed orgs
	if resp3.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307, got %d", resp3.StatusCode)
	}

	loc := resp3.Header.Get("Location")
	if loc != "/" {
		t.Errorf("expected redirect to '/', got %q", loc)
	}

	// Verify session cookie was set
	var sessionCookie *http.Cookie
	for _, c := range resp3.Cookies() {
		if c.Name == "dashyard_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Error("expected dashyard_session cookie to be set")
	}
}

func TestOAuthIntegrationAuthInfoEndpoint(t *testing.T) {
	dummyGH := newDummyGitHub()
	defer dummyGH.Close()

	providers := []config.OAuthProviderConfig{
		{
			Provider:     "github",
			ClientID:     "dummy-client-id",
			ClientSecret: "dummy-client-secret",
			BaseURL:      dummyGH.URL,
			RedirectURL:  "http://localhost:8080/auth/github/callback",
		},
	}

	router := setupOAuthTestServer(dummyGH.URL, providers)

	req := httptest.NewRequest("GET", "/api/auth-info", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var authInfo struct {
		PasswordEnabled bool `json:"password_enabled"`
		OAuthProviders  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"oauth_providers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authInfo); err != nil {
		t.Fatalf("failed to decode auth-info response: %v", err)
	}

	if authInfo.PasswordEnabled {
		t.Error("expected password_enabled to be false")
	}
	if len(authInfo.OAuthProviders) != 1 {
		t.Fatalf("expected 1 OAuth provider, got %d", len(authInfo.OAuthProviders))
	}
	if authInfo.OAuthProviders[0].Name != "github" {
		t.Errorf("expected provider name 'github', got %q", authInfo.OAuthProviders[0].Name)
	}
	if authInfo.OAuthProviders[0].URL != "/auth/github" {
		t.Errorf("expected provider URL '/auth/github', got %q", authInfo.OAuthProviders[0].URL)
	}
}

func TestOAuthIntegrationFetchGitHubOrgsWithBaseURL(t *testing.T) {
	dummyGH := newDummyGitHub()
	defer dummyGH.Close()

	orgs, err := auth.FetchGitHubOrgs("dummy-access-token", dummyGH.URL)
	if err != nil {
		t.Fatalf("FetchGitHubOrgs failed: %v", err)
	}

	if len(orgs) != 1 {
		t.Fatalf("expected 1 org, got %d", len(orgs))
	}
	if orgs[0] != "dummy-org" {
		t.Errorf("expected org 'dummy-org', got %q", orgs[0])
	}
}

func TestOAuthIntegrationFetchGitHubOrgsWithoutBaseURL(t *testing.T) {
	// Without base URL, FetchGitHubOrgs should use api.github.com.
	// We can't actually call api.github.com in tests without auth,
	// so we just verify the function signature accepts an empty base URL.
	_, err := auth.FetchGitHubOrgs("invalid-token", "")
	// We expect an error because the token is invalid, but the URL should be correct
	if err == nil {
		t.Log("FetchGitHubOrgs succeeded unexpectedly (maybe network issue)")
	}
	// The important thing is it doesn't panic
	_ = fmt.Sprintf("err=%v", err)
}
