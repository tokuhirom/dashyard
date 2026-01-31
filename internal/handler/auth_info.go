package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/config"
)

// AuthInfoResponse describes available authentication methods.
type AuthInfoResponse struct {
	PasswordEnabled bool   `json:"password_enabled"`
	OAuthEnabled    bool   `json:"oauth_enabled"`
	OAuthProvider   string `json:"oauth_provider,omitempty"`
	OAuthLoginURL   string `json:"oauth_login_url,omitempty"`
}

// AuthInfoHandler handles GET /api/auth-info requests.
type AuthInfoHandler struct {
	users    []config.User
	oauthCfg *config.OAuthConfig
}

// NewAuthInfoHandler creates a new AuthInfoHandler.
func NewAuthInfoHandler(users []config.User, oauthCfg *config.OAuthConfig) *AuthInfoHandler {
	return &AuthInfoHandler{
		users:    users,
		oauthCfg: oauthCfg,
	}
}

// Handle returns the available authentication methods.
func (h *AuthInfoHandler) Handle(c *gin.Context) {
	resp := AuthInfoResponse{
		PasswordEnabled: len(h.users) > 0,
		OAuthEnabled:    h.oauthCfg != nil,
	}

	if h.oauthCfg != nil {
		resp.OAuthProvider = h.oauthCfg.Provider
		resp.OAuthLoginURL = "/auth/login"
	}

	c.JSON(http.StatusOK, resp)
}
