// dummygithub is a fake GitHub OAuth server for local development and testing.
// It implements the minimal set of GitHub OAuth/API endpoints that Dashyard
// needs, allowing end-to-end testing without a real GitHub or GHE instance.
// Similar in spirit to cmd/dummyprom.
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	port := "5555"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("GET /login/oauth/authorize", handleAuthorize)
	http.HandleFunc("POST /login/oauth/access_token", handleAccessToken)
	http.HandleFunc("GET /api/v3/user", handleUser)
	http.HandleFunc("GET /api/v3/user/emails", handleEmails)
	http.HandleFunc("GET /api/v3/user/orgs", handleOrgs)

	slog.Info("dummy github server starting", "port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

// handleAuthorize shows a simple login form. When submitted, it redirects back
// to the caller's redirect_uri with a dummy authorization code.
func handleAuthorize(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")

	slog.Info("authorize", "redirect_uri", redirectURI, "state", state)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Dummy GitHub Login</title></head>
<body style="font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: #f6f8fa;">
  <div style="background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.12); text-align: center;">
    <h2>Dummy GitHub Login</h2>
    <p>Click the button to sign in as <strong>dummyuser</strong>.</p>
    <a href="%s?code=dummy-auth-code&state=%s"
       style="display: inline-block; padding: 0.6rem 1.5rem; background: #2da44e; color: white; text-decoration: none; border-radius: 6px; font-size: 1rem;">
      Sign in as dummyuser
    </a>
  </div>
</body>
</html>`, redirectURI, state)
}

// handleAccessToken exchanges the dummy code for a dummy access token.
func handleAccessToken(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	slog.Info("access_token", "code", code)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token": "dummy-access-token",
		"token_type":   "bearer",
		"scope":        "user:email,read:org",
	})
}

// handleUser returns a fixed GitHub user profile.
func handleUser(w http.ResponseWriter, r *http.Request) {
	slog.Info("user profile requested")

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"login":      "dummyuser",
		"id":         12345,
		"avatar_url": "https://avatars.githubusercontent.com/u/0?v=4",
		"name":       "Dummy User",
		"email":      "dummy@example.com",
	})
}

// handleEmails returns a fixed email list.
func handleEmails(w http.ResponseWriter, r *http.Request) {
	slog.Info("user emails requested")

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode([]map[string]interface{}{
		{
			"email":      "dummy@example.com",
			"primary":    true,
			"verified":   true,
			"visibility": "public",
		},
	})
}

// handleOrgs returns a fixed organization list.
func handleOrgs(w http.ResponseWriter, r *http.Request) {
	slog.Info("user orgs requested")

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode([]map[string]interface{}{
		{
			"login":       "dummy-org",
			"id":          100,
			"description": "A dummy organization for testing",
		},
	})
}
