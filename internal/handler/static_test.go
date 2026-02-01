package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

func testFS() fstest.MapFS {
	return fstest.MapFS{
		"index.html":           {Data: []byte("<html><body>SPA</body></html>")},
		"assets/style.css":     {Data: []byte("body { margin: 0; }")},
		"assets/app.js":        {Data: []byte("console.log('app')")},
		"assets/logo.svg":      {Data: []byte("<svg></svg>")},
		"favicon.ico":          {Data: []byte("icon-data")},
	}
}

func setupStaticRouter(fs fstest.MapFS) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := NewStaticHandler(fs)
	router.NoRoute(h.Handle)
	return router
}

func TestStaticHandlerServesExactFile(t *testing.T) {
	router := setupStaticRouter(testFS())

	req := httptest.NewRequest("GET", "/assets/style.css", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
	if body := resp.Body.String(); body != "body { margin: 0; }" {
		t.Errorf("expected CSS content, got %q", body)
	}
	ct := resp.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/css") {
		t.Errorf("expected Content-Type containing 'text/css', got %q", ct)
	}
}

func TestStaticHandlerServesJS(t *testing.T) {
	router := setupStaticRouter(testFS())

	req := httptest.NewRequest("GET", "/assets/app.js", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
	ct := resp.Header().Get("Content-Type")
	if !strings.Contains(ct, "javascript") {
		t.Errorf("expected Content-Type containing 'javascript', got %q", ct)
	}
}

func TestStaticHandlerSPAFallback(t *testing.T) {
	router := setupStaticRouter(testFS())

	paths := []string{
		"/dashboard/overview",
		"/d/some/nested/path",
		"/nonexistent",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest("GET", path, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Errorf("expected 200 for SPA fallback, got %d", resp.Code)
			}
			if body := resp.Body.String(); !strings.Contains(body, "SPA") {
				t.Errorf("expected index.html content, got %q", body)
			}
		})
	}
}

func TestStaticHandlerRootServesIndexHTML(t *testing.T) {
	router := setupStaticRouter(testFS())

	// Root path should serve index.html via the file server
	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
	if body := resp.Body.String(); !strings.Contains(body, "SPA") {
		t.Errorf("expected index.html content, got %q", body)
	}
	ct := resp.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Errorf("expected Content-Type containing 'text/html', got %q", ct)
	}
}

func TestStaticHandlerServesFavicon(t *testing.T) {
	router := setupStaticRouter(testFS())

	req := httptest.NewRequest("GET", "/favicon.ico", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
	if body := resp.Body.String(); body != "icon-data" {
		t.Errorf("expected favicon content, got %q", body)
	}
}
