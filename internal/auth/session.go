package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	cookieName    = "dashyard_session"
	sessionExpiry = 24 * time.Hour
)

// SessionPayload is the JSON structure stored in the session cookie.
type SessionPayload struct {
	UserID string `json:"user_id"`
	Exp    int64  `json:"exp"`
}

// SessionManager handles session cookie creation and validation.
type SessionManager struct {
	secret []byte
	secure bool // Set Secure flag on cookies (for HTTPS)
}

// NewSessionManager creates a new SessionManager with the given secret.
func NewSessionManager(secret string, secure bool) *SessionManager {
	return &SessionManager{
		secret: []byte(secret),
		secure: secure,
	}
}

// CreateSession sets a signed session cookie on the response.
func (sm *SessionManager) CreateSession(w http.ResponseWriter, userID string) error {
	payload := SessionPayload{
		UserID: userID,
		Exp:    time.Now().Add(sessionExpiry).Unix(),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling session payload: %w", err)
	}

	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)
	sig := sm.sign(payloadB64)
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	cookieValue := payloadB64 + "." + sigB64

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   sm.secure,
		MaxAge:   int(sessionExpiry.Seconds()),
	})

	return nil
}

// ValidateSession reads and validates the session cookie, returning the payload if valid.
func (sm *SessionManager) ValidateSession(r *http.Request) (*SessionPayload, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, fmt.Errorf("no session cookie")
	}

	parts := strings.SplitN(cookie.Value, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cookie format")
	}

	payloadB64 := parts[0]
	sigB64 := parts[1]

	sig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return nil, fmt.Errorf("invalid signature encoding")
	}

	expectedSig := sm.sign(payloadB64)
	if !hmac.Equal(sig, expectedSig) {
		return nil, fmt.Errorf("invalid signature")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, fmt.Errorf("invalid payload encoding")
	}

	var payload SessionPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	if time.Now().Unix() > payload.Exp {
		return nil, fmt.Errorf("session expired")
	}

	return &payload, nil
}

// ClearSession removes the session cookie.
func (sm *SessionManager) ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func (sm *SessionManager) sign(data string) []byte {
	mac := hmac.New(sha256.New, sm.secret)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}
