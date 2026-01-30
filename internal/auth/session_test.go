package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateAndValidateSession(t *testing.T) {
	sm := NewSessionManager("test-secret", false)

	// Create session
	w := httptest.NewRecorder()
	if err := sm.CreateSession(w, "admin"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Extract cookie from response
	resp := w.Result()
	cookies := resp.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != "dashyard_session" {
		t.Errorf("expected cookie name 'dashyard_session', got %q", cookie.Name)
	}
	if !cookie.HttpOnly {
		t.Error("expected HttpOnly flag")
	}
	if cookie.SameSite != http.SameSiteStrictMode {
		t.Error("expected SameSite=Strict")
	}

	// Validate session
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(cookie)
	payload, err := sm.ValidateSession(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payload.UserID != "admin" {
		t.Errorf("expected user_id 'admin', got %q", payload.UserID)
	}
	if payload.Exp <= time.Now().Unix() {
		t.Error("expected expiry to be in the future")
	}
}

func TestValidateSessionNoCookie(t *testing.T) {
	sm := NewSessionManager("test-secret", false)
	r := httptest.NewRequest("GET", "/", nil)
	_, err := sm.ValidateSession(r)
	if err == nil {
		t.Error("expected error for missing cookie")
	}
}

func TestValidateSessionInvalidSignature(t *testing.T) {
	sm1 := NewSessionManager("secret-1", false)
	sm2 := NewSessionManager("secret-2", false)

	w := httptest.NewRecorder()
	if err := sm1.CreateSession(w, "admin"); err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(w.Result().Cookies()[0])

	_, err := sm2.ValidateSession(r)
	if err == nil {
		t.Error("expected error for invalid signature")
	}
}

func TestValidateSessionTampered(t *testing.T) {
	sm := NewSessionManager("test-secret", false)
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "dashyard_session",
		Value: "tampered.value",
	})
	_, err := sm.ValidateSession(r)
	if err == nil {
		t.Error("expected error for tampered cookie")
	}
}

func TestClearSession(t *testing.T) {
	sm := NewSessionManager("test-secret", false)
	w := httptest.NewRecorder()
	sm.ClearSession(w)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].MaxAge != -1 {
		t.Errorf("expected MaxAge -1, got %d", cookies[0].MaxAge)
	}
}
