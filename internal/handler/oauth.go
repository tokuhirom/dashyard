package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/tokuhirom/dashyard/internal/auth"
	"github.com/tokuhirom/dashyard/internal/config"
)

// OAuthHandler handles OAuth authentication flow.
type OAuthHandler struct {
	providers []config.OAuthProviderConfig
	session   *auth.SessionManager
}

// NewOAuthHandler creates a new OAuthHandler.
func NewOAuthHandler(providers []config.OAuthProviderConfig, session *auth.SessionManager) *OAuthHandler {
	return &OAuthHandler{
		providers: providers,
		session:   session,
	}
}

// BeginAuth starts the OAuth flow by redirecting to the provider.
func (h *OAuthHandler) BeginAuth(c *gin.Context) {
	provider := c.Param("provider")

	// Set provider in query so gothic can find it
	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// Callback handles the OAuth callback from the provider.
func (h *OAuthHandler) Callback(c *gin.Context) {
	provider := c.Param("provider")

	// Set provider in query so gothic can find it
	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		slog.Error("OAuth callback failed", "error", err)
		c.Redirect(http.StatusTemporaryRedirect, "/?error=oauth_failed")
		return
	}

	// Find provider config for allowlist check
	providerCfg := auth.FindOAuthProvider(h.providers, provider)
	if providerCfg == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/?error=unknown_provider")
		return
	}

	allowed, err := auth.CheckUserAllowed(gothUser, *providerCfg)
	if err != nil {
		slog.Error("OAuth allowlist check failed", "error", err)
		c.Redirect(http.StatusTemporaryRedirect, "/?error=oauth_failed")
		return
	}
	if !allowed {
		c.Redirect(http.StatusTemporaryRedirect, "/?error=access_denied")
		return
	}

	// Use NickName as the user ID (GitHub username)
	userID := gothUser.NickName
	if userID == "" {
		userID = gothUser.Email
	}

	if err := h.session.CreateSession(c.Request, c.Writer, userID); err != nil {
		slog.Error("OAuth session creation failed", "error", err)
		c.Redirect(http.StatusTemporaryRedirect, "/?error=session_failed")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// Logout clears the session and redirects to the login page.
func (h *OAuthHandler) Logout(c *gin.Context) {
	if err := h.session.ClearSession(c.Request, c.Writer); err != nil {
		slog.Error("logout failed", "error", err)
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
