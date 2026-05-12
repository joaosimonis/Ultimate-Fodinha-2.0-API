package auth

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const adminRole = "admin"

type UserReader interface {
	FindByUsername(ctx context.Context, username string) (User, error)
}

type Service struct {
	repo   UserReader
	tokens *TokenManager
}

func NewService(repo UserReader, tokens *TokenManager) *Service {
	return &Service{repo: repo, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		return LoginResponse{}, ErrInvalidCredentials
	}

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return LoginResponse{}, ErrInvalidCredentials
		}
		return LoginResponse{}, err
	}

	if user.Role != adminRole {
		return LoginResponse{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}

	accessToken, err := s.tokens.Generate(user)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Authenticated: true,
		AccessToken:   accessToken,
		TokenType:     tokenTypeBearer,
		ExpiresIn:     s.tokens.ExpiresInSeconds(),
		User: AdminUser{
			Username: user.Username,
			Role:     user.Role,
		},
	}, nil
}
