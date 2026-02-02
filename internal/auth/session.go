package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	sessionName   = "dashyard_session"
	sessionUserID = "user_id"
	sessionMaxAge = 86400 // 24 hours
)

// SessionManager handles session creation and validation using gorilla/sessions.
type SessionManager struct {
	store *sessions.CookieStore
}

// NewSessionManager creates a new SessionManager with the given secret.
func NewSessionManager(secret string, secure bool) *SessionManager {
	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionMaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	}
	return &SessionManager{store: store}
}

// CreateSession saves a session with the given user ID.
func (sm *SessionManager) CreateSession(r *http.Request, w http.ResponseWriter, userID string) error {
	session, err := sm.store.Get(r, sessionName)
	if err != nil {
		// If the existing cookie is corrupt, create a fresh session
		session, err = sm.store.New(r, sessionName)
		if err != nil {
			return fmt.Errorf("creating session: %w", err)
		}
	}
	session.Values[sessionUserID] = userID
	return session.Save(r, w)
}

// ValidateSession reads the session and returns the user ID if valid.
func (sm *SessionManager) ValidateSession(r *http.Request) (string, error) {
	session, err := sm.store.Get(r, sessionName)
	if err != nil {
		return "", fmt.Errorf("invalid session: %w", err)
	}
	if session.IsNew {
		return "", fmt.Errorf("no session")
	}
	userID, ok := session.Values[sessionUserID].(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("no user_id in session")
	}
	return userID, nil
}

// ClearSession removes the session.
func (sm *SessionManager) ClearSession(r *http.Request, w http.ResponseWriter) error {
	session, err := sm.store.Get(r, sessionName)
	if err != nil {
		return fmt.Errorf("getting session: %w", err)
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

// ExpireCookie writes a Set-Cookie header that expires the session cookie.
// Unlike ClearSession, this works even when the existing cookie is corrupt
// or was created by a different instance with a different secret.
func (sm *SessionManager) ExpireCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   sm.store.Options.Secure,
	})
}

// Store returns the underlying CookieStore (used by gothic).
func (sm *SessionManager) Store() *sessions.CookieStore {
	return sm.store
}
