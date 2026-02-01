package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/markbates/goth"
	"github.com/tokuhirom/dashyard/internal/config"
)

func TestCheckUserAllowedNoRestrictions(t *testing.T) {
	user := goth.User{NickName: "anyone"}
	cfg := config.OAuthProviderConfig{Provider: "github"}

	allowed, err := CheckUserAllowed(user, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Error("expected user to be allowed when no restrictions")
	}
}

func TestCheckUserAllowedByUsername(t *testing.T) {
	user := goth.User{NickName: "user1"}
	cfg := config.OAuthProviderConfig{
		Provider:     "github",
		AllowedUsers: []string{"user1", "user2"},
	}

	allowed, err := CheckUserAllowed(user, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Error("expected user to be allowed by username")
	}
}

func TestCheckUserDeniedByUsername(t *testing.T) {
	user := goth.User{NickName: "hacker"}
	cfg := config.OAuthProviderConfig{
		Provider:     "github",
		AllowedUsers: []string{"user1", "user2"},
	}

	allowed, err := CheckUserAllowed(user, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Error("expected user to be denied")
	}
}

func TestInitGothProviders(t *testing.T) {
	providers := []config.OAuthProviderConfig{
		{
			Provider:     "github",
			ClientID:     "test-id",
			ClientSecret: "test-secret",
			RedirectURL:  "http://localhost:8080/auth/github/callback",
		},
	}

	// Should not panic
	InitGothProviders(providers)

	// Verify provider was registered
	p, err := goth.GetProvider("github")
	if err != nil {
		t.Fatalf("expected github provider to be registered: %v", err)
	}
	if p.Name() != "github" {
		t.Errorf("expected provider name 'github', got %q", p.Name())
	}
}

func TestInitGothProvidersWithBaseURL(t *testing.T) {
	providers := []config.OAuthProviderConfig{
		{
			Provider:     "github",
			ClientID:     "test-id",
			ClientSecret: "test-secret",
			RedirectURL:  "http://localhost:8080/auth/github/callback",
			BaseURL:      "https://ghe.example.com",
		},
	}

	// Should not panic
	InitGothProviders(providers)

	// Verify provider was registered
	p, err := goth.GetProvider("github")
	if err != nil {
		t.Fatalf("expected github provider to be registered: %v", err)
	}
	if p.Name() != "github" {
		t.Errorf("expected provider name 'github', got %q", p.Name())
	}
}

func TestFetchGitHubOrgs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"login":"org1"},{"login":"org2"}]`))
	}))
	defer server.Close()

	// FetchGitHubOrgs with baseURL uses baseURL + "/api/v3/user/orgs"
	orgs, err := FetchGitHubOrgs("test-token", server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 2 {
		t.Fatalf("expected 2 orgs, got %d", len(orgs))
	}
	if orgs[0] != "org1" || orgs[1] != "org2" {
		t.Errorf("expected [org1, org2], got %v", orgs)
	}
}

func TestFetchGitHubOrgsDefaultURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"login":"myorg"}]`))
	}))
	defer server.Close()

	// Can't easily test the default github.com URL without mocking,
	// but we can test the baseURL path with our server.
	orgs, err := FetchGitHubOrgs("token", server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 1 || orgs[0] != "myorg" {
		t.Errorf("expected [myorg], got %v", orgs)
	}
}

func TestFetchGitHubOrgsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	_, err := FetchGitHubOrgs("bad-token", server.URL)
	if err == nil {
		t.Error("expected error for non-200 status")
	}
}

func TestFetchGitHubOrgsInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	_, err := FetchGitHubOrgs("token", server.URL)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestCheckUserAllowedByOrg(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"login":"allowed-org"},{"login":"other-org"}]`))
	}))
	defer server.Close()

	user := goth.User{NickName: "someuser", AccessToken: "test-token"}
	cfg := config.OAuthProviderConfig{
		Provider:    "github",
		AllowedOrgs: []string{"allowed-org"},
		BaseURL:     server.URL,
	}

	allowed, err := CheckUserAllowed(user, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Error("expected user to be allowed by org membership")
	}
}

func TestCheckUserDeniedByOrg(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"login":"other-org"}]`))
	}))
	defer server.Close()

	user := goth.User{NickName: "someuser", AccessToken: "test-token"}
	cfg := config.OAuthProviderConfig{
		Provider:    "github",
		AllowedOrgs: []string{"required-org"},
		BaseURL:     server.URL,
	}

	allowed, err := CheckUserAllowed(user, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Error("expected user to be denied when not in required org")
	}
}

func TestCheckUserAllowedOrgFetchError(t *testing.T) {
	user := goth.User{NickName: "someuser", AccessToken: "test-token"}
	cfg := config.OAuthProviderConfig{
		Provider:    "github",
		AllowedOrgs: []string{"some-org"},
		BaseURL:     "http://localhost:1", // unreachable
	}

	_, err := CheckUserAllowed(user, cfg)
	if err == nil {
		t.Error("expected error when org fetch fails")
	}
}

func TestFindOAuthProvider(t *testing.T) {
	providers := []config.OAuthProviderConfig{
		{Provider: "github", ClientID: "gh-id"},
	}

	p := FindOAuthProvider(providers, "github")
	if p == nil {
		t.Fatal("expected to find github provider")
	}
	if p.ClientID != "gh-id" {
		t.Errorf("expected client_id 'gh-id', got %q", p.ClientID)
	}

	p = FindOAuthProvider(providers, "google")
	if p != nil {
		t.Error("expected nil for unknown provider")
	}
}
