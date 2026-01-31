package auth

import (
	"context"
	"fmt"

	"github.com/tokuhirom/dashyard/internal/config"
	"golang.org/x/oauth2"
)

// OAuthUserInfo holds user information retrieved from an OAuth provider.
type OAuthUserInfo struct {
	ID       string
	Username string
	Email    string
	Orgs     []string // GitHub-specific
}

// OAuthProvider defines the interface for OAuth/OIDC providers.
type OAuthProvider interface {
	// AuthCodeURL returns the URL to redirect the user to for authentication.
	AuthCodeURL(state string) string
	// Exchange exchanges an authorization code for a token.
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	// UserInfo fetches the authenticated user's info using the given token.
	UserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error)
	// Name returns the display name of the provider (e.g. "GitHub", "Google").
	Name() string
}

// NewOAuthProvider creates an OAuthProvider from the given config.
func NewOAuthProvider(cfg *config.OAuthConfig) (OAuthProvider, error) {
	switch cfg.Provider {
	case "github":
		return NewGitHubProvider(cfg), nil
	case "google":
		return NewOIDCProvider(cfg, "https://accounts.google.com"), nil
	case "oidc":
		return NewOIDCProvider(cfg, cfg.IssuerURL), nil
	default:
		return nil, fmt.Errorf("unsupported oauth provider: %q", cfg.Provider)
	}
}

// IsUserAllowed checks whether the user is allowed by the configured allowlists.
// If no allowlists are configured, all users are allowed.
func IsUserAllowed(cfg *config.OAuthConfig, info *OAuthUserInfo) bool {
	if len(cfg.AllowedUsers) == 0 && len(cfg.AllowedOrgs) == 0 {
		return true
	}

	for _, u := range cfg.AllowedUsers {
		if u == info.Username {
			return true
		}
	}

	for _, allowedOrg := range cfg.AllowedOrgs {
		for _, userOrg := range info.Orgs {
			if allowedOrg == userOrg {
				return true
			}
		}
	}

	return false
}
