package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAndValidateSession(t *testing.T) {
	sm := NewSessionManager("test-secret-that-is-32bytes!!", false)

	// Create session
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	if err := sm.CreateSession(r, w, "admin"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Extract cookie from response
	resp := w.Result()
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected at least 1 cookie")
	}

	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "dashyard_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("expected dashyard_session cookie")
	}
	if !sessionCookie.HttpOnly {
		t.Error("expected HttpOnly flag")
	}

	// Validate session
	r2 := httptest.NewRequest("GET", "/", nil)
	r2.AddCookie(sessionCookie)
	userID, err := sm.ValidateSession(r2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != "admin" {
		t.Errorf("expected user_id 'admin', got %q", userID)
	}
}

func TestValidateSessionNoCookie(t *testing.T) {
	sm := NewSessionManager("test-secret-that-is-32bytes!!", false)
	r := httptest.NewRequest("GET", "/", nil)
	_, err := sm.ValidateSession(r)
	if err == nil {
		t.Error("expected error for missing cookie")
	}
}

func TestValidateSessionInvalidSignature(t *testing.T) {
	sm1 := NewSessionManager("secret-1-that-is-at-least-32!!", false)
	sm2 := NewSessionManager("secret-2-that-is-at-least-32!!", false)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	if err := sm1.CreateSession(r, w, "admin"); err != nil {
		t.Fatal(err)
	}

	r2 := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		r2.AddCookie(c)
	}

	_, err := sm2.ValidateSession(r2)
	if err == nil {
		t.Error("expected error for invalid signature")
	}
}

func TestValidateSessionTampered(t *testing.T) {
	sm := NewSessionManager("test-secret-that-is-32bytes!!", false)
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

func TestCreateSessionOverwritesExisting(t *testing.T) {
	sm := NewSessionManager("test-secret-that-is-32bytes!!", false)

	// Create session for user1
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	if err := sm.CreateSession(r, w, "user1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create a new session for user2 using the same cookie
	r2 := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		r2.AddCookie(c)
	}
	w2 := httptest.NewRecorder()
	if err := sm.CreateSession(r2, w2, "user2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate the session now has user2
	r3 := httptest.NewRequest("GET", "/", nil)
	for _, c := range w2.Result().Cookies() {
		r3.AddCookie(c)
	}
	userID, err := sm.ValidateSession(r3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != "user2" {
		t.Errorf("expected user_id 'user2', got %q", userID)
	}
}

func TestValidateSessionEmptyUserID(t *testing.T) {
	sm := NewSessionManager("test-secret-that-is-32bytes!!", false)

	// Create session with empty user ID
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	if err := sm.CreateSession(r, w, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r2 := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		r2.AddCookie(c)
	}
	_, err := sm.ValidateSession(r2)
	if err == nil {
		t.Error("expected error for empty user_id in session")
	}
}

func TestNewSessionManagerSecureFlag(t *testing.T) {
	sm := NewSessionManager("test-secret-that-is-32bytes!!", true)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	if err := sm.CreateSession(r, w, "admin"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, c := range w.Result().Cookies() {
		if c.Name == "dashyard_session" {
			if !c.Secure {
				t.Error("expected Secure flag on cookie")
			}
			return
		}
	}
	t.Fatal("session cookie not found")
}

func TestStoreReturnsUnderlyingStore(t *testing.T) {
	sm := NewSessionManager("test-secret-that-is-32bytes!!", false)
	store := sm.Store()
	if store == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestClearSession(t *testing.T) {
	sm := NewSessionManager("test-secret-that-is-32bytes!!", false)

	// Create a session first
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	if err := sm.CreateSession(r, w, "admin"); err != nil {
		t.Fatal(err)
	}

	// Clear it
	r2 := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		r2.AddCookie(c)
	}
	w2 := httptest.NewRecorder()
	if err := sm.ClearSession(r2, w2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := w2.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected at least 1 cookie")
	}
	found := false
	for _, c := range cookies {
		if c.Name == "dashyard_session" && c.MaxAge < 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected session cookie with MaxAge < 0")
	}
}
