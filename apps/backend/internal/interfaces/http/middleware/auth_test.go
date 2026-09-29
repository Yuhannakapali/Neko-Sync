package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nekosync/internal/infrastructure/auth"

	"github.com/labstack/echo/v4"
)

func run(t *testing.T, header string) (int, any) {
	t.Helper()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if header != "" {
		req.Header.Set("Authorization", header)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	verifier := auth.NewJWTManager("secret", time.Hour)
	h := AuthMiddleware(verifier)(func(c echo.Context) error { return c.NoContent(http.StatusOK) })
	_ = h(c)
	return rec.Code, c.Get("user_id")
}

func TestAuthMiddlewareRejectsForgedToken(t *testing.T) {
	for _, header := range []string{"", "Token abc", "Bearer ", "Bearer anything-at-all"} {
		if code, _ := run(t, header); code != http.StatusUnauthorized {
			t.Errorf("header %q: got %d, want 401", header, code)
		}
	}
}

func TestAuthMiddlewareAcceptsValidToken(t *testing.T) {
	tok, err := auth.NewJWTManager("secret", time.Hour).Issue("user-42")
	if err != nil {
		t.Fatal(err)
	}
	code, userID := run(t, "Bearer "+tok)
	if code != http.StatusOK || userID != "user-42" {
		t.Fatalf("got %d with user_id %v", code, userID)
	}
}
