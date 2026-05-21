package matches

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

const matchColumns = "id, mode, status, room_code, max_score, lives, live_state, created_at, finished_at, cancelled_at"

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) FindAll(ctx context.Context, status string) ([]Match, error) {
	query := "SELECT " + matchColumns + " FROM matches ORDER BY created_at DESC"
	args := []any{}

	if status != "" {
		query = "SELECT " + matchColumns + " FROM matches WHERE status = $1 ORDER BY created_at DESC"
		args = append(args, status)
	}

	rows, err := r.pool.Query(ctx, query, args...)
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
	match, err := r.findOne(ctx, "SELECT "+matchColumns+" FROM matches WHERE room_code = $1", roomCode)
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
	match, err := r.findOne(ctx, "SELECT "+matchColumns+" FROM matches WHERE id = $1", id)
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

func (r *Repository) Create(ctx context.Context, req CreateMatchRequest, roomCode string, hostTokenHash string) (Match, error) {
	match, err := scanMatch(r.pool.QueryRow(ctx, `
		INSERT INTO matches (mode, room_code, host_token_hash, max_score, lives)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+matchColumns+`
	`, req.Mode, roomCode, hostTokenHash, req.MaxScore, req.Lives))
	if err != nil {
		if isUniqueViolation(err) {
			return Match{}, ErrRoomCodeTaken
		}
		return Match{}, err
	}

	return match, nil
}

func (r *Repository) AddPlayers(ctx context.Context, id string, req AddMatchPlayersRequest) (Match, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Match{}, err
	}
	defer tx.Rollback(ctx)

	for _, player := range req.Players {
		if _, err := tx.Exec(ctx, `
			INSERT INTO match_players (match_id, player_id, team_name)
			VALUES ($1, $2, $3)
		`, id, player.PlayerID, normalizedTeamName(player.TeamName)); err != nil {
			if isUniqueViolation(err) {
				return Match{}, ErrDuplicatePlayer
			}
			return Match{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Match{}, err
	}

	return r.FindByID(ctx, id)
}

func (r *Repository) RemovePlayer(ctx context.Context, id string, playerID string) (Match, error) {
	result, err := r.pool.Exec(ctx, `
		DELETE FROM match_players
		WHERE match_id = $1 AND player_id = $2
	`, id, playerID)
	if err != nil {
		return Match{}, err
	}

	if result.RowsAffected() == 0 {
		return Match{}, ErrNotFound
	}

	return r.FindByID(ctx, id)
}

func (r *Repository) Finish(ctx context.Context, id string, results []FinishMatchPlayerResult) (Match, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Match{}, err
	}
	defer tx.Rollback(ctx)

	for _, player := range results {
		if _, err := tx.Exec(ctx, `
			UPDATE match_players
			SET final_score = $1,
				remaining_lives = $2,
				placement = $3,
				is_winner = $4
			WHERE match_id = $5 AND player_id = $6
		`, player.FinalScore, player.RemainingLives, player.Placement, player.IsWinner, id, player.PlayerID); err != nil {
			return Match{}, err
		}
	}

	match, err := scanMatch(tx.QueryRow(ctx, `
		UPDATE matches
		SET status = 'finished', finished_at = NOW()
		WHERE id = $1
		RETURNING `+matchColumns+`
	`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Match{}, ErrNotFound
		}
		return Match{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Match{}, err
	}

	players, err := r.findPlayersByMatchID(ctx, match.ID)
	if err != nil {
		return Match{}, err
	}
	match.Players = players

	return match, nil
}

func (r *Repository) CancelByRoomCode(ctx context.Context, roomCode string) (Match, error) {
	match, err := scanMatch(r.pool.QueryRow(ctx, `
		UPDATE matches
		SET status = 'cancelled',
			cancelled_at = COALESCE(cancelled_at, NOW())
		WHERE LOWER(room_code) = LOWER($1)
		RETURNING `+matchColumns+`
	`, roomCode))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Match{}, ErrNotFound
		}
		return Match{}, err
	}

	players, err := r.findPlayersByMatchID(ctx, match.ID)
	if err != nil {
		return Match{}, err
	}
	match.Players = players

	return match, nil
}

func (r *Repository) ExistsByRoomCode(ctx context.Context, roomCode string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM matches WHERE LOWER(room_code) = LOWER($1))
	`, roomCode).Scan(&exists)
	return exists, err
}

func (r *Repository) FindHostTokenHashByRoomCode(ctx context.Context, roomCode string) (string, error) {
	var hostTokenHash sql.NullString
	err := r.pool.QueryRow(ctx, `
		SELECT host_token_hash
		FROM matches
		WHERE LOWER(room_code) = LOWER($1)
	`, roomCode).Scan(&hostTokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}

	if !hostTokenHash.Valid || strings.TrimSpace(hostTokenHash.String) == "" {
		return "", ErrInvalidHostToken
	}

	return hostTokenHash.String, nil
}

func (r *Repository) UpdateLiveStateByRoomCode(ctx context.Context, roomCode string, liveState []byte) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE matches
		SET live_state = $1
		WHERE LOWER(room_code) = LOWER($2)
	`, liveState, roomCode)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
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
		SELECT id, match_id, player_id, team_name, final_score, remaining_lives, placement, is_winner
		FROM match_players
		WHERE match_id = $1
		ORDER BY team_name NULLS LAST, id
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
	var liveState []byte
	var finishedAt sql.NullTime
	var cancelledAt sql.NullTime

	err := scanner.Scan(
		&match.ID,
		&match.Mode,
		&match.Status,
		&match.RoomCode,
		&maxScore,
		&lives,
		&liveState,
		&match.CreatedAt,
		&finishedAt,
		&cancelledAt,
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

	if len(liveState) > 0 && string(liveState) != "null" {
		match.LiveState = liveState
	}

	if finishedAt.Valid {
		match.FinishedAt = &finishedAt.Time
	}

	if cancelledAt.Valid {
		match.CancelledAt = &cancelledAt.Time
	}

	return match, nil
}

func scanMatchPlayer(scanner matchScanner) (MatchPlayer, error) {
	var player MatchPlayer
	var teamName sql.NullString
	var finalScore sql.NullInt64
	var remainingLives sql.NullInt64
	var placement sql.NullInt64

	err := scanner.Scan(
		&player.ID,
		&player.MatchID,
		&player.PlayerID,
		&teamName,
		&finalScore,
		&remainingLives,
		&placement,
		&player.IsWinner,
	)
	if err != nil {
		return MatchPlayer{}, err
	}

	if teamName.Valid {
		player.TeamName = teamName.String
	}

	if finalScore.Valid {
		value := int(finalScore.Int64)
		player.FinalScore = &value
	}

	if remainingLives.Valid {
		value := int(remainingLives.Int64)
		player.RemainingLives = &value
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

func normalizedTeamName(teamName string) any {
	teamName = strings.TrimSpace(teamName)
	if teamName == "" {
		return nil
	}

	return teamName
}
