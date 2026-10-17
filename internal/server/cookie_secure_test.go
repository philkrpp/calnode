package server

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/calnode/calnode/internal/config"
	"github.com/calnode/calnode/internal/db"
)

// An https instance with only email/password login (no Google or Microsoft OAuth)
// must still mark its session cookie Secure.
func TestSessionCookieSecureWithoutOAuth(t *testing.T) {
	database, err := db.Open("sqlite://:memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := db.Migrate(database); err != nil {
		t.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx, cancel := context.WithCancel(t.Context())
	h, drain := BuildHandler(ctx, &config.Config{
		BaseURL: "https://calnode.localhost", DataDir: t.TempDir(), CookieSecure: true,
	}, database, logger)
	defer func() { cancel(); drain() }()

	rec := httptest.NewRecorder()
	h.Claim(rec, httptest.NewRequest(http.MethodPost, "/v1/auth/claim", strings.NewReader(`{"name":"Host","email":"host@example.com","password":"correct-horse-battery","timezone":"UTC"}`)))
	if rec.Code >= 300 {
		t.Fatalf("claim: %d %s", rec.Code, rec.Body)
	}
	var found bool
	for _, c := range rec.Result().Cookies() {
		if c.Name != "calnode_session" {
			continue
		}
		found = true
		if !c.Secure {
			t.Error("calnode_session cookie is not Secure on an https instance without OAuth")
		}
	}
	if !found {
		t.Fatal("claim did not set calnode_session")
	}
}
