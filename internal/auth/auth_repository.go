package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (User, error) {
	username = strings.TrimSpace(username)

	var user User
	err := r.pool.QueryRow(ctx, `
		SELECT id, username, password_hash, role
		FROM users
		WHERE LOWER(username) = LOWER($1)
	`, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, nil
}
