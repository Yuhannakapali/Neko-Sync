package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestIssueAndVerify(t *testing.T) {
	m := NewJWTManager("secret", time.Hour)
	tok, err := m.Issue("user-1")
	if err != nil {
		t.Fatal(err)
	}
	id, err := m.Verify(tok)
	if err != nil || id != "user-1" {
		t.Fatalf("got %q, %v", id, err)
	}
}

func TestVerifyRejects(t *testing.T) {
	m := NewJWTManager("secret", time.Hour)

	other, _ := NewJWTManager("other-secret", time.Hour).Issue("user-1")

	expiredMgr := NewJWTManager("secret", time.Hour)
	expiredMgr.now = func() time.Time { return time.Now().Add(-2 * time.Hour) }
	expired, _ := expiredMgr.Issue("user-1")

	none, _ := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
		Subject:   "user-1",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString(jwt.UnsafeAllowNoneSignatureType)

	noExpiry, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Subject: "user-1"}).SignedString([]byte("secret"))

	noSubject, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString([]byte("secret"))

	cases := map[string]string{
		"garbage":      "not-a-jwt",
		"wrong secret": other,
		"expired":      expired,
		"alg none":     none,
		"no expiry":    noExpiry,
		"no subject":   noSubject,
	}
	for name, tok := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := m.Verify(tok); err != ErrInvalidToken {
				t.Fatalf("expected ErrInvalidToken, got %v", err)
			}
		})
	}
}
