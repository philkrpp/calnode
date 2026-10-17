package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/calnode/calnode/frontend"
)

func TestSecurityHeaders_httpsInstance(t *testing.T) {
	h := SecurityHeaders(true)(frontend.Handler())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	for header, want := range map[string]string{
		"Strict-Transport-Security": "max-age=31536000",
		"X-Content-Type-Options":    "nosniff",
		"Referrer-Policy":           "same-origin",
	} {
		if got := rec.Header().Get(header); got != want {
			t.Errorf("%s = %q; want %q", header, got, want)
		}
	}
}

// Plain-http dev instances must not pin HSTS on localhost.
func TestSecurityHeaders_noHSTSOverHTTP(t *testing.T) {
	h := SecurityHeaders(false)(frontend.Handler())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if got := rec.Header().Get("Strict-Transport-Security"); got != "" {
		t.Errorf("Strict-Transport-Security = %q; want it absent", got)
	}
}

// A handler that sets its own Referrer-Policy keeps it.
func TestSecurityHeaders_handlerOverrides(t *testing.T) {
	h := SecurityHeaders(true)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "no-referrer")
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if got := rec.Header().Get("Referrer-Policy"); got != "no-referrer" {
		t.Errorf("Referrer-Policy = %q; want no-referrer", got)
	}
}
