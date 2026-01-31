package acl

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewEmpty(t *testing.T) {
	al, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if al != nil {
		t.Error("expected nil for empty allow list")
	}

	al, err = New([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if al != nil {
		t.Error("expected nil for empty allow list")
	}
}

func TestNewCIDR(t *testing.T) {
	al, err := New([]string{"192.168.1.0/24", "10.0.0.0/8"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if al == nil {
		t.Fatal("expected non-nil allow list")
	}
	if len(al.networks) != 2 {
		t.Errorf("expected 2 networks, got %d", len(al.networks))
	}
}

func TestNewSingleIP(t *testing.T) {
	al, err := New([]string{"192.168.1.100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if al == nil {
		t.Fatal("expected non-nil allow list")
	}
	if al.networks[0].String() != "192.168.1.100/32" {
		t.Errorf("expected 192.168.1.100/32, got %s", al.networks[0].String())
	}
}

func TestNewIPv6(t *testing.T) {
	al, err := New([]string{"::1", "fd00::/8"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if al == nil {
		t.Fatal("expected non-nil allow list")
	}
	if len(al.networks) != 2 {
		t.Errorf("expected 2 networks, got %d", len(al.networks))
	}
}

func TestNewInvalid(t *testing.T) {
	_, err := New([]string{"not-an-ip"})
	if err == nil {
		t.Error("expected error for invalid entry")
	}
}

func TestContainsMatch(t *testing.T) {
	al, err := New([]string{"192.168.1.0/24"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !al.Contains("192.168.1.50") {
		t.Error("expected 192.168.1.50 to be contained")
	}
}

func TestContainsNoMatch(t *testing.T) {
	al, err := New([]string{"192.168.1.0/24"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if al.Contains("10.0.0.1") {
		t.Error("expected 10.0.0.1 to not be contained")
	}
}

func TestContainsInvalidIP(t *testing.T) {
	al, err := New([]string{"192.168.1.0/24"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if al.Contains("invalid") {
		t.Error("expected invalid IP to not be contained")
	}
}

func TestMiddlewareNilAllowAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Middleware(nil))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMiddlewareAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	al, err := New([]string{"192.168.1.0/24"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := gin.New()
	r.Use(Middleware(al))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.50:12345"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMiddlewareDenied(t *testing.T) {
	gin.SetMode(gin.TestMode)
	al, err := New([]string{"192.168.1.0/24"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := gin.New()
	r.Use(Middleware(al))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}
