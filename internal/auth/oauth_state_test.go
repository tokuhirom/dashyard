package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOAuthStateGenerateAndValidate(t *testing.T) {
	m := NewOAuthStateManager("test-secret", false)

	w := httptest.NewRecorder()
	state, err := m.Generate(w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state == "" {
		t.Fatal("expected non-empty state")
	}

	// Check cookie was set
	cookies := w.Result().Cookies()
	var stateCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == oauthStateCookieName {
			stateCookie = c
			break
		}
	}
	if stateCookie == nil {
		t.Fatal("expected state cookie to be set")
	}
	if stateCookie.SameSite != http.SameSiteLaxMode {
		t.Error("expected SameSite=Lax")
	}
	if !stateCookie.HttpOnly {
		t.Error("expected HttpOnly")
	}

	// Validate
	r := httptest.NewRequest("GET", "/auth/callback?state="+state, nil)
	r.AddCookie(stateCookie)
	w2 := httptest.NewRecorder()
	if err := m.Validate(w2, r, state); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	// Cookie should be cleared after validation
	clearCookies := w2.Result().Cookies()
	for _, c := range clearCookies {
		if c.Name == oauthStateCookieName && c.MaxAge == -1 {
			return // OK
		}
	}
	t.Error("expected state cookie to be cleared after validation")
}

func TestOAuthStateValidateMismatch(t *testing.T) {
	m := NewOAuthStateManager("test-secret", false)

	w := httptest.NewRecorder()
	state, err := m.Generate(w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(w.Result().Cookies()[0])
	w2 := httptest.NewRecorder()

	err = m.Validate(w2, r, state+"tampered")
	if err == nil {
		t.Error("expected error for mismatched state")
	}
}

func TestOAuthStateValidateNoCookie(t *testing.T) {
	m := NewOAuthStateManager("test-secret", false)
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	err := m.Validate(w, r, "some-state")
	if err == nil {
		t.Error("expected error for missing cookie")
	}
}

func TestOAuthStateValidateWrongSecret(t *testing.T) {
	m1 := NewOAuthStateManager("secret-1", false)
	m2 := NewOAuthStateManager("secret-2", false)

	w := httptest.NewRecorder()
	state, err := m1.Generate(w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(w.Result().Cookies()[0])
	w2 := httptest.NewRecorder()

	err = m2.Validate(w2, r, state)
	if err == nil {
		t.Error("expected error for wrong secret")
	}
}
