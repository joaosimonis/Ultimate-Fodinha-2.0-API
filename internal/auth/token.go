package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenTypeBearer = "Bearer"

type TokenManager struct {
	secret           string
	expiresInSeconds int
	now              func() time.Time
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenManager(secret string, expiresInSeconds int) *TokenManager {
	return &TokenManager{
		secret:           secret,
		expiresInSeconds: expiresInSeconds,
		now:              time.Now,
	}
}

func (m *TokenManager) ExpiresInSeconds() int {
	return m.expiresInSeconds
}

func (m *TokenManager) Generate(user User) (string, error) {
	now := m.now()
	claims := Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.expiresInSeconds) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *TokenManager) Validate(tokenString string) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(m.secret), nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(m.now),
	)

	if err != nil || token == nil || !token.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return Claims{}, ErrExpiredToken
		}

		return Claims{}, ErrInvalidToken
	}

	if strings.TrimSpace(claims.Subject) == "" || strings.TrimSpace(claims.Username) == "" || claims.Role != adminRole {
		return Claims{}, ErrInvalidToken
	}

	return claims, nil
}

func BearerToken(authorization string) (string, error) {
	tokenType, token, found := strings.Cut(strings.TrimSpace(authorization), " ")
	if !found || !strings.EqualFold(tokenType, tokenTypeBearer) || strings.TrimSpace(token) == "" {
		return "", ErrInvalidToken
	}

	return strings.TrimSpace(token), nil
}

func TokenErrorMessage(err error) string {
	if errors.Is(err, ErrExpiredToken) {
		return ErrExpiredToken.Error()
	}

	if errors.Is(err, ErrInvalidToken) {
		return ErrInvalidToken.Error()
	}

	return fmt.Sprint(err)
}
