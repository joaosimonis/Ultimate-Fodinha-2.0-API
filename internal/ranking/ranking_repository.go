package ranking

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) FindFodinha(ctx context.Context, limit int) ([]Entry, error) {
	rows, err := r.pool.Query(ctx, `
		WITH stats AS (
			SELECT
				p.id AS player_id,
				p.name,
				p.avatar_key,
				COUNT(*)::INT AS games_played,
				COUNT(*) FILTER (WHERE mp.placement = 1)::INT AS gold_medals,
				COUNT(*) FILTER (WHERE mp.placement = 2)::INT AS silver_medals,
				COUNT(*) FILTER (WHERE mp.placement = 3)::INT AS bronze_medals,
				COUNT(*) FILTER (WHERE mp.placement BETWEEN 1 AND 3)::INT AS podiums
			FROM match_players mp
			INNER JOIN matches m ON m.id = mp.match_id
			INNER JOIN players p ON p.id = mp.player_id
			WHERE m.mode = 'fodinha'
				AND m.status = 'finished'
			GROUP BY p.id, p.name, p.avatar_key
		)
		SELECT
			ROW_NUMBER() OVER (
				ORDER BY gold_medals DESC,
					silver_medals DESC,
					bronze_medals DESC,
					podiums DESC,
					games_played DESC,
					name ASC
			)::INT AS position,
			player_id,
			name,
			avatar_key,
			games_played,
			gold_medals,
			silver_medals,
			bronze_medals,
			podiums
		FROM stats
		ORDER BY position ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		var goldMedals int
		var silverMedals int
		var bronzeMedals int
		var podiums int

		if err := rows.Scan(
			&entry.Position,
			&entry.PlayerID,
			&entry.Name,
			&entry.AvatarKey,
			&entry.GamesPlayed,
			&goldMedals,
			&silverMedals,
			&bronzeMedals,
			&podiums,
		); err != nil {
			return nil, err
		}

		entry.GoldMedals = intPtr(goldMedals)
		entry.SilverMedals = intPtr(silverMedals)
		entry.BronzeMedals = intPtr(bronzeMedals)
		entry.Podiums = intPtr(podiums)
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if entries == nil {
		return []Entry{}, nil
	}

	return entries, nil
}

func (r *Repository) FindTruco(ctx context.Context, teamSize int, limit int) ([]Entry, error) {
	rows, err := r.pool.Query(ctx, `
		WITH team_counts AS (
			SELECT
				m.id AS match_id,
				mp.team_name,
				COUNT(*)::INT AS team_players
			FROM matches m
			INNER JOIN match_players mp ON mp.match_id = m.id
			WHERE m.mode = 'truco'
				AND m.status = 'finished'
			GROUP BY m.id, mp.team_name
		),
		truco_matches AS (
			SELECT match_id
			FROM team_counts
			GROUP BY match_id
			HAVING COUNT(*) = 2
				AND MIN(team_players) = $1
				AND MAX(team_players) = $1
		),
		stats AS (
			SELECT
				p.id AS player_id,
				p.name,
				p.avatar_key,
				COUNT(*)::INT AS games_played,
				COUNT(*) FILTER (WHERE mp.is_winner)::INT AS wins,
				COUNT(*) FILTER (WHERE NOT mp.is_winner)::INT AS losses,
				ROUND((COUNT(*) FILTER (WHERE mp.is_winner)::NUMERIC * 100 / COUNT(*)), 2)::FLOAT8 AS win_rate
			FROM match_players mp
			INNER JOIN truco_matches tm ON tm.match_id = mp.match_id
			INNER JOIN players p ON p.id = mp.player_id
			GROUP BY p.id, p.name, p.avatar_key
		)
		SELECT
			ROW_NUMBER() OVER (
				ORDER BY wins DESC,
					win_rate DESC,
					games_played DESC,
					name ASC
			)::INT AS position,
			player_id,
			name,
			avatar_key,
			games_played,
			wins,
			losses,
			win_rate
		FROM stats
		ORDER BY position ASC
		LIMIT $2
	`, teamSize, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		var wins int
		var losses int
		var winRate float64

		if err := rows.Scan(
			&entry.Position,
			&entry.PlayerID,
			&entry.Name,
			&entry.AvatarKey,
			&entry.GamesPlayed,
			&wins,
			&losses,
			&winRate,
		); err != nil {
			return nil, err
		}

		entry.Wins = intPtr(wins)
		entry.Losses = intPtr(losses)
		entry.WinRate = floatPtr(winRate)
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if entries == nil {
		return []Entry{}, nil
	}

	return entries, nil
}

func intPtr(value int) *int {
	return &value
}

func floatPtr(value float64) *float64 {
	return &value
}
