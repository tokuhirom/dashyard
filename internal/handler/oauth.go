package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/auth"
	"github.com/tokuhirom/dashyard/internal/config"
)

// OAuthHandler handles OAuth login, callback, and logout routes.
type OAuthHandler struct {
	provider     auth.OAuthProvider
	stateManager *auth.OAuthStateManager
	session      *auth.SessionManager
	oauthCfg     *config.OAuthConfig
}

// NewOAuthHandler creates a new OAuthHandler.
func NewOAuthHandler(provider auth.OAuthProvider, stateManager *auth.OAuthStateManager, session *auth.SessionManager, oauthCfg *config.OAuthConfig) *OAuthHandler {
	return &OAuthHandler{
		provider:     provider,
		stateManager: stateManager,
		session:      session,
		oauthCfg:     oauthCfg,
	}
}

// Login redirects the user to the OAuth provider's authorization page.
func (h *OAuthHandler) Login(c *gin.Context) {
	state, err := h.stateManager.Generate(c.Writer)
	if err != nil {
		log.Printf("oauth: failed to generate state: %v", err)
		c.String(http.StatusInternalServerError, "Failed to initiate login")
		return
	}

	url := h.provider.AuthCodeURL(state)
	if url == "" {
		log.Printf("oauth: failed to generate auth URL")
		c.String(http.StatusInternalServerError, "Failed to initiate login")
		return
	}

	c.Redirect(http.StatusFound, url)
}

// Callback handles the OAuth callback after the user authorizes.
func (h *OAuthHandler) Callback(c *gin.Context) {
	// Validate state
	state := c.Query("state")
	if err := h.stateManager.Validate(c.Writer, c.Request, state); err != nil {
		log.Printf("oauth: state validation failed: %v", err)
		c.String(http.StatusBadRequest, "Invalid or expired OAuth state")
		return
	}

	// Check for error from provider
	if errParam := c.Query("error"); errParam != "" {
		desc := c.Query("error_description")
		log.Printf("oauth: provider returned error: %s: %s", errParam, desc)
		c.String(http.StatusBadRequest, "OAuth error: %s", desc)
		return
	}

	// Exchange code for token
	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Missing authorization code")
		return
	}

	token, err := h.provider.Exchange(c.Request.Context(), code)
	if err != nil {
		log.Printf("oauth: token exchange failed: %v", err)
		c.String(http.StatusInternalServerError, "Failed to exchange authorization code")
		return
	}

	// Fetch user info
	info, err := h.provider.UserInfo(c.Request.Context(), token)
	if err != nil {
		log.Printf("oauth: userinfo fetch failed: %v", err)
		c.String(http.StatusInternalServerError, "Failed to fetch user information")
		return
	}

	// Check allowlist
	if !auth.IsUserAllowed(h.oauthCfg, info) {
		log.Printf("oauth: user %q not allowed", info.Username)
		c.String(http.StatusForbidden, "Access denied")
		return
	}

	// Create session
	userID := info.Username
	if userID == "" {
		userID = info.ID
	}

	if err := h.session.CreateSession(c.Writer, userID); err != nil {
		log.Printf("oauth: session creation failed: %v", err)
		c.String(http.StatusInternalServerError, "Failed to create session")
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// Logout clears the session and redirects to the home page.
func (h *OAuthHandler) Logout(c *gin.Context) {
	h.session.ClearSession(c.Writer)
	c.Redirect(http.StatusFound, "/")
}
