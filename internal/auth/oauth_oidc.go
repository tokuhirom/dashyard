package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/tokuhirom/dashyard/internal/config"
	"golang.org/x/oauth2"
)

// oidcDiscovery holds the OIDC discovery document endpoints.
type oidcDiscovery struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
}

// OIDCProvider implements OAuthProvider for generic OIDC providers.
type OIDCProvider struct {
	cfg       *config.OAuthConfig
	issuerURL string
	name      string

	// Lazy-loaded discovery
	once      sync.Once
	discovery *oidcDiscovery
	discErr   error

	// Overridable HTTP client for testing
	httpClient *http.Client
}

// NewOIDCProvider creates a new OIDC provider.
func NewOIDCProvider(cfg *config.OAuthConfig, issuerURL string) *OIDCProvider {
	name := "OIDC"
	if cfg.Provider == "google" {
		name = "Google"
	}
	return &OIDCProvider{
		cfg:        cfg,
		issuerURL:  strings.TrimRight(issuerURL, "/"),
		name:       name,
		httpClient: http.DefaultClient,
	}
}

func (p *OIDCProvider) discover() (*oidcDiscovery, error) {
	p.once.Do(func() {
		url := p.issuerURL + "/.well-known/openid-configuration"
		resp, err := p.httpClient.Get(url)
		if err != nil {
			p.discErr = fmt.Errorf("fetching OIDC discovery: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			p.discErr = fmt.Errorf("OIDC discovery returned %d: %s", resp.StatusCode, body)
			return
		}

		var disc oidcDiscovery
		if err := json.NewDecoder(resp.Body).Decode(&disc); err != nil {
			p.discErr = fmt.Errorf("decoding OIDC discovery: %w", err)
			return
		}
		p.discovery = &disc
	})
	return p.discovery, p.discErr
}

func (p *OIDCProvider) oauthConfig() (*oauth2.Config, error) {
	disc, err := p.discover()
	if err != nil {
		return nil, err
	}

	scopes := p.cfg.Scopes
	if len(scopes) == 0 {
		scopes = []string{"openid", "profile", "email"}
	}

	return &oauth2.Config{
		ClientID:     p.cfg.ClientID,
		ClientSecret: p.cfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  disc.AuthorizationEndpoint,
			TokenURL: disc.TokenEndpoint,
		},
		RedirectURL: p.cfg.RedirectURL,
		Scopes:      scopes,
	}, nil
}

func (p *OIDCProvider) AuthCodeURL(state string) string {
	cfg, err := p.oauthConfig()
	if err != nil {
		// Fallback: return empty string if discovery fails
		return ""
	}
	return cfg.AuthCodeURL(state)
}

func (p *OIDCProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	cfg, err := p.oauthConfig()
	if err != nil {
		return nil, err
	}
	return cfg.Exchange(ctx, code)
}

func (p *OIDCProvider) UserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error) {
	disc, err := p.discover()
	if err != nil {
		return nil, err
	}

	cfg, err := p.oauthConfig()
	if err != nil {
		return nil, err
	}

	client := cfg.Client(ctx, token)
	resp, err := client.Get(disc.UserinfoEndpoint)
	if err != nil {
		return nil, fmt.Errorf("fetching OIDC userinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OIDC userinfo returned %d: %s", resp.StatusCode, body)
	}

	var result struct {
		Sub               string `json:"sub"`
		PreferredUsername string `json:"preferred_username"`
		Email             string `json:"email"`
		Name              string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding OIDC userinfo: %w", err)
	}

	username := result.PreferredUsername
	if username == "" {
		username = result.Email
	}
	if username == "" {
		username = result.Name
	}

	return &OAuthUserInfo{
		ID:       result.Sub,
		Username: username,
		Email:    result.Email,
	}, nil
}

func (p *OIDCProvider) Name() string {
	return p.name
}
