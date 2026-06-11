package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datewu/gtea"
)

func TestMyRoutesGzip(t *testing.T) {
	cfg := gtea.DefaultConfig()
	app := gtea.NewApp(context.Background(), cfg)
	h := New(app)

	// Test GET /my/logs with Accept-Encoding: gzip
	req := httptest.NewRequest(http.MethodGet, "/my/logs", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	
	t.Logf("Response status: %d", w.Code)
	t.Logf("Content-Encoding header: %q", w.Header().Get("Content-Encoding"))
	
	// Test GET /my/profile with Accept-Encoding: gzip
	req2 := httptest.NewRequest(http.MethodGet, "/my/profile", nil)
	req2.Header.Set("Accept-Encoding", "gzip")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, req2)
	
	t.Logf("Profile Response status: %d", w2.Code)
	t.Logf("Profile Content-Encoding header: %q", w2.Header().Get("Content-Encoding"))
}
