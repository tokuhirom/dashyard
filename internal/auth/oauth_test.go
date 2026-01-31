package auth

import (
	"testing"

	"github.com/tokuhirom/dashyard/internal/config"
)

func TestIsUserAllowed(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.OAuthConfig
		info    *OAuthUserInfo
		allowed bool
	}{
		{
			name:    "no allowlist allows everyone",
			cfg:     &config.OAuthConfig{},
			info:    &OAuthUserInfo{Username: "anyone"},
			allowed: true,
		},
		{
			name:    "allowed user",
			cfg:     &config.OAuthConfig{AllowedUsers: []string{"alice", "bob"}},
			info:    &OAuthUserInfo{Username: "alice"},
			allowed: true,
		},
		{
			name:    "disallowed user",
			cfg:     &config.OAuthConfig{AllowedUsers: []string{"alice"}},
			info:    &OAuthUserInfo{Username: "eve"},
			allowed: false,
		},
		{
			name:    "allowed org",
			cfg:     &config.OAuthConfig{AllowedOrgs: []string{"my-org"}},
			info:    &OAuthUserInfo{Username: "eve", Orgs: []string{"my-org", "other-org"}},
			allowed: true,
		},
		{
			name:    "disallowed org",
			cfg:     &config.OAuthConfig{AllowedOrgs: []string{"my-org"}},
			info:    &OAuthUserInfo{Username: "eve", Orgs: []string{"other-org"}},
			allowed: false,
		},
		{
			name:    "allowed by user when orgs also configured",
			cfg:     &config.OAuthConfig{AllowedUsers: []string{"alice"}, AllowedOrgs: []string{"my-org"}},
			info:    &OAuthUserInfo{Username: "alice", Orgs: []string{}},
			allowed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUserAllowed(tt.cfg, tt.info)
			if got != tt.allowed {
				t.Errorf("IsUserAllowed() = %v, want %v", got, tt.allowed)
			}
		})
	}
}

func TestNewOAuthProviderGitHub(t *testing.T) {
	cfg := &config.OAuthConfig{
		Provider:     "github",
		ClientID:     "id",
		ClientSecret: "secret",
		RedirectURL:  "http://localhost/callback",
	}
	p, err := NewOAuthProvider(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "GitHub" {
		t.Errorf("expected name 'GitHub', got %q", p.Name())
	}
}

func TestNewOAuthProviderOIDC(t *testing.T) {
	cfg := &config.OAuthConfig{
		Provider:     "oidc",
		ClientID:     "id",
		ClientSecret: "secret",
		RedirectURL:  "http://localhost/callback",
		IssuerURL:    "https://example.com",
	}
	p, err := NewOAuthProvider(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "OIDC" {
		t.Errorf("expected name 'OIDC', got %q", p.Name())
	}
}

func TestNewOAuthProviderGoogle(t *testing.T) {
	cfg := &config.OAuthConfig{
		Provider:     "google",
		ClientID:     "id",
		ClientSecret: "secret",
		RedirectURL:  "http://localhost/callback",
	}
	p, err := NewOAuthProvider(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "Google" {
		t.Errorf("expected name 'Google', got %q", p.Name())
	}
}

func TestNewOAuthProviderUnsupported(t *testing.T) {
	cfg := &config.OAuthConfig{
		Provider: "unsupported",
	}
	_, err := NewOAuthProvider(cfg)
	if err == nil {
		t.Error("expected error for unsupported provider")
	}
}
