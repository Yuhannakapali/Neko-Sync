package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

// JWTManager issues and verifies HS256 access tokens whose subject is the user ID.
type JWTManager struct {
	secret []byte
	expiry time.Duration
	now    func() time.Time
}

func NewJWTManager(secret string, expiry time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), expiry: expiry, now: time.Now}
}

func (m *JWTManager) Issue(userID string) (string, error) {
	now := m.now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.expiry)),
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

func (m *JWTManager) Verify(token string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(token, claims,
		func(*jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(m.now),
	)
	if err != nil || claims.Subject == "" {
		return "", ErrInvalidToken
	}
	return claims.Subject, nil
}
