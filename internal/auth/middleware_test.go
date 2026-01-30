package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthMiddlewareValid(t *testing.T) {
	sm := NewSessionManager("test-secret", false)

	router := gin.New()
	router.Use(AuthMiddleware(sm))
	router.GET("/test", func(c *gin.Context) {
		userID := GetUserID(c)
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	// Create a session cookie
	w := httptest.NewRecorder()
	if err := sm.CreateSession(w, "admin"); err != nil {
		t.Fatal(err)
	}
	cookie := w.Result().Cookies()[0]

	// Make authenticated request
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(cookie)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["user_id"] != "admin" {
		t.Errorf("expected user_id 'admin', got %q", body["user_id"])
	}
}

func TestAuthMiddlewareNoCookie(t *testing.T) {
	sm := NewSessionManager("test-secret", false)

	router := gin.New()
	router.Use(AuthMiddleware(sm))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestAuthMiddlewareInvalidCookie(t *testing.T) {
	sm := NewSessionManager("test-secret", false)

	router := gin.New()
	router.Use(AuthMiddleware(sm))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "dashyard_session", Value: "invalid"})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}
