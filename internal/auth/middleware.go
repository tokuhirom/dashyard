package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const userIDKey = "user_id"

// AuthMiddleware returns a Gin middleware that requires a valid session.
// It sets the user_id in the Gin context on success.
func AuthMiddleware(sm *SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := sm.ValidateSession(c.Request)
		if err != nil {
			// Clear the invalid/corrupt session cookie so re-login works cleanly
			sm.ExpireCookie(c.Writer)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
		c.Set(userIDKey, userID)
		c.Next()
	}
}

// GetUserID retrieves the authenticated user ID from the Gin context.
func GetUserID(c *gin.Context) string {
	v, _ := c.Get(userIDKey)
	s, _ := v.(string)
	return s
}
