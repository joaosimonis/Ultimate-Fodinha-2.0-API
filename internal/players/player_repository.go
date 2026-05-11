package players

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) FindAll(ctx context.Context) ([]Player, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, avatar_key, title, lore, ultimate, created_at
		FROM players
		ORDER BY avatar_key ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		var ultimateJSON []byte
		if err := rows.Scan(&p.ID, &p.Name, &p.AvatarKey, &p.Title, &p.Lore, &ultimateJSON, &p.CreatedAt); err != nil {
			return nil, err
		}
		if ultimateJSON != nil {
			if err := json.Unmarshal(ultimateJSON, &p.Ultimate); err != nil {
				return nil, err
			}
		}
		players = append(players, p)
	}
	return players, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (Player, error) {
	var p Player
	var ultimateJSON []byte

	err := r.pool.QueryRow(ctx, `
		SELECT id, name, avatar_key, title, lore, ultimate, created_at
		FROM players
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &p.AvatarKey, &p.Title, &p.Lore, &ultimateJSON, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Player{}, ErrNotFound
		}
		return Player{}, err
	}

	if ultimateJSON != nil {
		if err := json.Unmarshal(ultimateJSON, &p.Ultimate); err != nil {
			return Player{}, err
		}
	}

	return p, nil
}

func (r *Repository) Create(ctx context.Context, req CreatePlayerRequest) (Player, error) {
	ultimateJSON, err := json.Marshal(req.Ultimate)
	if err != nil {
		return Player{}, err
	}

	var p Player
	var ultimateBytes []byte
	err = r.pool.QueryRow(ctx, `
		INSERT INTO players (name, avatar_key, title, lore, ultimate)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, avatar_key, title, lore, ultimate, created_at
	`, req.Name, req.AvatarKey, req.Title, req.Lore, ultimateJSON).
		Scan(&p.ID, &p.Name, &p.AvatarKey, &p.Title, &p.Lore, &ultimateBytes, &p.CreatedAt)
	if err != nil {
		return Player{}, err
	}

	if ultimateBytes != nil {
		if err := json.Unmarshal(ultimateBytes, &p.Ultimate); err != nil {
			return Player{}, err
		}
	}

	return p, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, `
		DELETE FROM players WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM players WHERE LOWER(name) = LOWER($1))
	`, name).Scan(&exists)
	return exists, err
}
