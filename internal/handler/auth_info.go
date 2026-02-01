package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/config"
)

// AuthInfoResponse is the JSON response for GET /api/auth-info.
type AuthInfoResponse struct {
	PasswordEnabled bool                `json:"password_enabled"`
	OAuthProviders  []OAuthProviderInfo `json:"oauth_providers"`
}

// OAuthProviderInfo describes an available OAuth provider for the frontend.
type OAuthProviderInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// AuthInfoHandler handles GET /api/auth-info.
type AuthInfoHandler struct {
	users     []config.User
	providers []config.OAuthProviderConfig
}

// NewAuthInfoHandler creates a new AuthInfoHandler.
func NewAuthInfoHandler(users []config.User, providers []config.OAuthProviderConfig) *AuthInfoHandler {
	return &AuthInfoHandler{
		users:     users,
		providers: providers,
	}
}

// Handle returns the authentication methods available.
func (h *AuthInfoHandler) Handle(c *gin.Context) {
	resp := AuthInfoResponse{
		PasswordEnabled: len(h.users) > 0,
		OAuthProviders:  make([]OAuthProviderInfo, 0, len(h.providers)),
	}

	for _, p := range h.providers {
		resp.OAuthProviders = append(resp.OAuthProviders, OAuthProviderInfo{
			Name: p.Provider,
			URL:  "/auth/" + p.Provider,
		})
	}

	c.JSON(http.StatusOK, resp)
}
