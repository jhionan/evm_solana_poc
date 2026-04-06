package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Role identifies an access role for a user.
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// Claims is the JWT payload embedded in every token.
type Claims struct {
	Wallet string `json:"wallet"`
	Role   Role   `json:"role"`
	jwt.RegisteredClaims
}

// JWTService issues and verifies signed JWTs.
type JWTService struct {
	secret []byte
}

// NewJWTService constructs a JWTService using the provided HMAC secret.
func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: []byte(secret)}
}

// Sign creates a signed HS256 token for the given claims valid for duration.
// IssuedAt and ExpiresAt are set by this method; callers should not set them.
func (s *JWTService) Sign(claims Claims, duration time.Duration) (string, error) {
	now := time.Now()
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(duration))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("auth: sign token: %w", err)
	}
	return signed, nil
}

// Verify parses and validates a token string, returning its Claims.
// Returns an error if the token is expired, has been tampered with, or uses
// an unexpected signing method.
func (s *JWTService) Verify(tokenStr string) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("auth: unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("auth: verify token: %w", err)
	}
	if !token.Valid {
		return Claims{}, errors.New("auth: token is not valid")
	}
	return claims, nil
}
