package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/auth"
	"github.com/tokuhirom/dashyard/internal/config"
)

type loginRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler handles POST /api/login requests.
type LoginHandler struct {
	users   []config.User
	session *auth.SessionManager
}

// NewLoginHandler creates a new LoginHandler.
func NewLoginHandler(users []config.User, session *auth.SessionManager) *LoginHandler {
	return &LoginHandler{
		users:   users,
		session: session,
	}
}

// Handle processes a login request.
func (h *LoginHandler) Handle(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Find user
	var user *config.User
	for i := range h.users {
		if h.users[i].ID == req.UserID {
			user = &h.users[i]
			break
		}
	}

	if user == nil || !auth.VerifyPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := h.session.CreateSession(c.Request, c.Writer, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": user.ID})
}
