package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/dashyard/internal/auth"
	"github.com/tokuhirom/dashyard/internal/config"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func generateTestHash(password string) string {
	c := crypt.SHA512.New()
	hash, err := c.Generate([]byte(password), nil)
	if err != nil {
		panic(err)
	}
	return hash
}

func TestLoginSuccess(t *testing.T) {
	users := []config.User{
		{ID: "admin", PasswordHash: generateTestHash("password123")},
	}
	sm := auth.NewSessionManager("test-secret", false)
	handler := NewLoginHandler(users, sm)

	router := gin.New()
	router.POST("/api/login", handler.Handle)

	body := `{"user_id":"admin","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["user_id"] != "admin" {
		t.Errorf("expected user_id 'admin', got %q", result["user_id"])
	}

	// Check session cookie was set
	cookies := resp.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "dashyard_session" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected session cookie to be set")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	users := []config.User{
		{ID: "admin", PasswordHash: generateTestHash("password123")},
	}
	sm := auth.NewSessionManager("test-secret", false)
	handler := NewLoginHandler(users, sm)

	router := gin.New()
	router.POST("/api/login", handler.Handle)

	body := `{"user_id":"admin","password":"wrongpassword"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestLoginUnknownUser(t *testing.T) {
	users := []config.User{
		{ID: "admin", PasswordHash: generateTestHash("password123")},
	}
	sm := auth.NewSessionManager("test-secret", false)
	handler := NewLoginHandler(users, sm)

	router := gin.New()
	router.POST("/api/login", handler.Handle)

	body := `{"user_id":"unknown","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestLoginBadRequest(t *testing.T) {
	sm := auth.NewSessionManager("test-secret", false)
	handler := NewLoginHandler(nil, sm)

	router := gin.New()
	router.POST("/api/login", handler.Handle)

	body := `{"invalid":"json"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}
