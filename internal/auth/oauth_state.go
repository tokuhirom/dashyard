package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net/http"
	"time"
)

const (
	oauthStateCookieName = "dashyard_oauth_state"
	stateExpiry          = 10 * time.Minute
	nonceLen             = 16
)

// OAuthStateManager generates and validates HMAC-signed OAuth state parameters.
type OAuthStateManager struct {
	secret []byte
	secure bool
}

// NewOAuthStateManager creates a new OAuthStateManager.
func NewOAuthStateManager(secret string, secure bool) *OAuthStateManager {
	return &OAuthStateManager{
		secret: []byte(secret),
		secure: secure,
	}
}

// Generate creates a state parameter and sets it as a cookie.
// The state is: base64(nonce + timestamp + hmac(nonce + timestamp)).
func (m *OAuthStateManager) Generate(w http.ResponseWriter) (string, error) {
	nonce := make([]byte, nonceLen)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("generating nonce: %w", err)
	}

	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts, uint64(time.Now().Unix()))

	payload := append(nonce, ts...)
	mac := hmac.New(sha256.New, m.secret)
	mac.Write(payload)
	sig := mac.Sum(nil)

	raw := append(payload, sig...)
	state := base64.RawURLEncoding.EncodeToString(raw)

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   m.secure,
		MaxAge:   int(stateExpiry.Seconds()),
	})

	return state, nil
}

// Validate checks the state parameter against the cookie and verifies the HMAC and timestamp.
// It clears the cookie after validation.
func (m *OAuthStateManager) Validate(w http.ResponseWriter, r *http.Request, state string) error {
	cookie, err := r.Cookie(oauthStateCookieName)
	if err != nil {
		return fmt.Errorf("missing state cookie")
	}

	// Clear the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	if state != cookie.Value {
		return fmt.Errorf("state mismatch")
	}

	raw, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		return fmt.Errorf("invalid state encoding")
	}

	// nonce (nonceLen) + timestamp (8) + HMAC-SHA256 (32)
	expectedLen := nonceLen + 8 + sha256.Size
	if len(raw) != expectedLen {
		return fmt.Errorf("invalid state length")
	}

	payload := raw[:nonceLen+8]
	sig := raw[nonceLen+8:]

	mac := hmac.New(sha256.New, m.secret)
	mac.Write(payload)
	expectedSig := mac.Sum(nil)

	if !hmac.Equal(sig, expectedSig) {
		return fmt.Errorf("invalid state signature")
	}

	ts := binary.BigEndian.Uint64(raw[nonceLen : nonceLen+8])
	created := time.Unix(int64(ts), 0)
	if time.Since(created) > stateExpiry {
		return fmt.Errorf("state expired")
	}

	return nil
}
