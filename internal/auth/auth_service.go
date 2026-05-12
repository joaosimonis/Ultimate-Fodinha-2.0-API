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
	repo UserReader
}

func NewService(repo UserReader) *Service {
	return &Service{repo: repo}
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

	return LoginResponse{
		Authenticated: true,
		User: AdminUser{
			Username: user.Username,
			Role:     user.Role,
		},
	}, nil
}
