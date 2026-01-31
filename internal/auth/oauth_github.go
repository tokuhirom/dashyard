package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tokuhirom/dashyard/internal/config"
	"golang.org/x/oauth2"
)

const (
	githubAuthURL    = "https://github.com/login/oauth/authorize"
	githubTokenURL   = "https://github.com/login/oauth/access_token"
	githubUserURL    = "https://api.github.com/user"
	githubOrgsURL    = "https://api.github.com/user/orgs"
)

// GitHubProvider implements OAuthProvider for GitHub.
type GitHubProvider struct {
	config   *oauth2.Config
	oauthCfg *config.OAuthConfig
	// Overridable for testing
	userURL string
	orgsURL string
}

// NewGitHubProvider creates a new GitHubProvider.
func NewGitHubProvider(cfg *config.OAuthConfig) *GitHubProvider {
	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []string{"read:user", "read:org"}
	}

	return &GitHubProvider{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  githubAuthURL,
				TokenURL: githubTokenURL,
			},
			RedirectURL: cfg.RedirectURL,
			Scopes:      scopes,
		},
		oauthCfg: cfg,
		userURL:  githubUserURL,
		orgsURL:  githubOrgsURL,
	}
}

func (p *GitHubProvider) AuthCodeURL(state string) string {
	return p.config.AuthCodeURL(state)
}

func (p *GitHubProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.config.Exchange(ctx, code)
}

func (p *GitHubProvider) UserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error) {
	client := p.config.Client(ctx, token)

	// Fetch user info
	userResp, err := client.Get(p.userURL)
	if err != nil {
		return nil, fmt.Errorf("fetching github user: %w", err)
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(userResp.Body)
		return nil, fmt.Errorf("github user API returned %d: %s", userResp.StatusCode, body)
	}

	var userResult struct {
		Login string `json:"login"`
		ID    int64  `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&userResult); err != nil {
		return nil, fmt.Errorf("decoding github user: %w", err)
	}

	info := &OAuthUserInfo{
		ID:       fmt.Sprintf("%d", userResult.ID),
		Username: userResult.Login,
		Email:    userResult.Email,
	}

	// Fetch orgs if allowed_orgs is configured
	if len(p.oauthCfg.AllowedOrgs) > 0 {
		orgs, err := p.fetchOrgs(ctx, client)
		if err != nil {
			return nil, err
		}
		info.Orgs = orgs
	}

	return info, nil
}

func (p *GitHubProvider) fetchOrgs(ctx context.Context, client *http.Client) ([]string, error) {
	resp, err := client.Get(p.orgsURL)
	if err != nil {
		return nil, fmt.Errorf("fetching github orgs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github orgs API returned %d: %s", resp.StatusCode, body)
	}

	var orgsResult []struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&orgsResult); err != nil {
		return nil, fmt.Errorf("decoding github orgs: %w", err)
	}

	orgs := make([]string, len(orgsResult))
	for i, o := range orgsResult {
		orgs[i] = o.Login
	}
	return orgs, nil
}

func (p *GitHubProvider) Name() string {
	return "GitHub"
}
