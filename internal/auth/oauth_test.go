package auth

import (
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
