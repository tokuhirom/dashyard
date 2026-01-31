package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}

	// Verify counter was incremented
	count := testutil.ToFloat64(HTTPRequestsTotal.WithLabelValues("GET", "/test", "200"))
	if count < 1 {
		t.Errorf("expected request counter >= 1, got %f", count)
	}

	// Verify histogram was observed
	histCount := testutil.CollectAndCount(HTTPRequestDuration)
	if histCount == 0 {
		t.Error("expected histogram to have observations")
	}
}
