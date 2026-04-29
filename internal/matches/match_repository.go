package matches

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) FindLive(ctx context.Context) ([]Match, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, mode, status, room_code, max_score, lives, created_at, finished_at
		FROM matches
		WHERE status = 'ongoing'
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		match, err := scanMatch(rows)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func (r *Repository) FindByRoomCode(ctx context.Context, roomCode string) (Match, error) {
	match, err := r.findOne(ctx, `
		SELECT id, mode, status, room_code, max_score, lives, created_at, finished_at
		FROM matches
		WHERE room_code = $1
	`, roomCode)
	if err != nil {
		return Match{}, err
	}

	players, err := r.findPlayersByMatchID(ctx, match.ID)
	if err != nil {
		return Match{}, err
	}
	match.Players = players

	return match, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (Match, error) {
	match, err := r.findOne(ctx, `
		SELECT id, mode, status, room_code, max_score, lives, created_at, finished_at
		FROM matches
		WHERE id = $1
	`, id)
	if err != nil {
		return Match{}, err
	}

	players, err := r.findPlayersByMatchID(ctx, match.ID)
	if err != nil {
		return Match{}, err
	}
	match.Players = players

	return match, nil
}

func (r *Repository) Create(ctx context.Context, req CreateMatchRequest) (Match, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Match{}, err
	}
	defer tx.Rollback(ctx)

	match, err := scanMatch(tx.QueryRow(ctx, `
		INSERT INTO matches (mode, room_code, max_score, lives)
		VALUES ($1, $2, $3, $4)
		RETURNING id, mode, status, room_code, max_score, lives, created_at, finished_at
	`, req.Mode, req.RoomCode, req.MaxScore, req.Lives))
	if err != nil {
		if isUniqueViolation(err) {
			return Match{}, ErrRoomCodeTaken
		}
		return Match{}, err
	}

	for _, player := range req.Players {
		matchPlayer, err := scanMatchPlayer(tx.QueryRow(ctx, `
			INSERT INTO match_players (match_id, player_id, team_index)
			VALUES ($1, $2, $3)
			RETURNING id, match_id, player_id, team_index, final_score, placement, is_winner
		`, match.ID, player.PlayerID, player.TeamIndex))
		if err != nil {
			return Match{}, err
		}
		match.Players = append(match.Players, matchPlayer)
	}

	if err := tx.Commit(ctx); err != nil {
		return Match{}, err
	}

	return match, nil
}

func (r *Repository) Finish(ctx context.Context, id string) (Match, error) {
	return r.findOne(ctx, `
		UPDATE matches
		SET status = 'finished', finished_at = NOW()
		WHERE id = $1
		RETURNING id, mode, status, room_code, max_score, lives, created_at, finished_at
	`, id)
}

func (r *Repository) ExistsByRoomCode(ctx context.Context, roomCode string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM matches WHERE LOWER(room_code) = LOWER($1))
	`, roomCode).Scan(&exists)
	return exists, err
}

func (r *Repository) findOne(ctx context.Context, query string, args ...any) (Match, error) {
	match, err := scanMatch(r.pool.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Match{}, ErrNotFound
		}
		return Match{}, err
	}

	return match, nil
}

func (r *Repository) findPlayersByMatchID(ctx context.Context, matchID string) ([]MatchPlayer, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, match_id, player_id, team_index, final_score, placement, is_winner
		FROM match_players
		WHERE match_id = $1
		ORDER BY team_index NULLS LAST, id
	`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []MatchPlayer
	for rows.Next() {
		player, err := scanMatchPlayer(rows)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

type matchScanner interface {
	Scan(dest ...any) error
}

func scanMatch(scanner matchScanner) (Match, error) {
	var match Match
	var maxScore sql.NullInt64
	var lives sql.NullInt64
	var finishedAt sql.NullTime

	err := scanner.Scan(
		&match.ID,
		&match.Mode,
		&match.Status,
		&match.RoomCode,
		&maxScore,
		&lives,
		&match.CreatedAt,
		&finishedAt,
	)
	if err != nil {
		return Match{}, err
	}

	if maxScore.Valid {
		value := int(maxScore.Int64)
		match.MaxScore = &value
	}

	if lives.Valid {
		value := int(lives.Int64)
		match.Lives = &value
	}

	if finishedAt.Valid {
		match.FinishedAt = &finishedAt.Time
	}

	return match, nil
}

func scanMatchPlayer(scanner matchScanner) (MatchPlayer, error) {
	var player MatchPlayer
	var teamIndex sql.NullInt64
	var placement sql.NullInt64

	err := scanner.Scan(
		&player.ID,
		&player.MatchID,
		&player.PlayerID,
		&teamIndex,
		&player.FinalScore,
		&placement,
		&player.IsWinner,
	)
	if err != nil {
		return MatchPlayer{}, err
	}

	if teamIndex.Valid {
		value := int(teamIndex.Int64)
		player.TeamIndex = &value
	}

	if placement.Valid {
		value := int(placement.Int64)
		player.Placement = &value
	}

	return player, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
}
