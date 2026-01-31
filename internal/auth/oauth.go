package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/tokuhirom/dashyard/internal/config"
)

// InitGothProviders configures goth with the OAuth providers from config.
func InitGothProviders(providers []config.OAuthProviderConfig) {
	var gothProviders []goth.Provider
	for _, p := range providers {
		switch p.Provider {
		case "github":
			scopes := p.Scopes
			if len(scopes) == 0 {
				scopes = []string{"user:email"}
			}
			gp := github.New(p.ClientID, p.ClientSecret, p.RedirectURL, scopes...)
			gothProviders = append(gothProviders, gp)
		}
	}
	goth.UseProviders(gothProviders...)
}

// CheckUserAllowed checks whether a goth user is permitted by the provider's allowlist.
// If no allowed_users and no allowed_orgs are configured, all authenticated users are allowed.
func CheckUserAllowed(user goth.User, providerCfg config.OAuthProviderConfig) (bool, error) {
	hasRestrictions := len(providerCfg.AllowedUsers) > 0 || len(providerCfg.AllowedOrgs) > 0

	if !hasRestrictions {
		return true, nil
	}

	// Check allowed users
	if slices.Contains(providerCfg.AllowedUsers, user.NickName) {
		return true, nil
	}

	// Check allowed orgs (GitHub-specific)
	if len(providerCfg.AllowedOrgs) > 0 && providerCfg.Provider == "github" {
		orgs, err := FetchGitHubOrgs(user.AccessToken)
		if err != nil {
			return false, fmt.Errorf("fetching GitHub orgs: %w", err)
		}
		for _, org := range orgs {
			if slices.Contains(providerCfg.AllowedOrgs, org) {
				return true, nil
			}
		}
	}

	return false, nil
}

// FetchGitHubOrgs retrieves the list of organization login names for the authenticated user.
func FetchGitHubOrgs(accessToken string) ([]string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/orgs", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var orgs []struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&orgs); err != nil {
		return nil, err
	}

	names := make([]string, len(orgs))
	for i, o := range orgs {
		names[i] = o.Login
	}
	return names, nil
}

// FindOAuthProvider finds a provider config by provider name.
func FindOAuthProvider(providers []config.OAuthProviderConfig, name string) *config.OAuthProviderConfig {
	for i := range providers {
		if providers[i].Provider == name {
			return &providers[i]
		}
	}
	return nil
}
