package ranking

import (
	"context"
	"encoding/json"

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
					avatar_key ASC
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
					avatar_key ASC
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

func (r *Repository) FindTrucoTeams(ctx context.Context, teamSize int, limit int) ([]TeamEntry, error) {
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
		team_results AS (
			SELECT
				m.id AS match_id,
				COALESCE(m.finished_at, m.created_at) AS played_at,
				mp.team_name,
				ARRAY_AGG(mp.player_id ORDER BY mp.player_id) AS player_ids,
				STRING_AGG(mp.player_id::TEXT, ':' ORDER BY mp.player_id) AS team_key,
				BOOL_OR(mp.is_winner) AS is_winner,
				MAX(mp.final_score)::INT AS final_score
			FROM matches m
			INNER JOIN truco_matches tm ON tm.match_id = m.id
			INNER JOIN match_players mp ON mp.match_id = m.id
			GROUP BY m.id, played_at, mp.team_name
			HAVING COUNT(*) = $1
		),
		team_matchups AS (
			SELECT
				team_results.team_key,
				team_results.player_ids,
				team_results.played_at,
				team_results.is_winner,
				team_results.final_score,
				opponents.final_score AS opponent_score
			FROM team_results
			INNER JOIN team_results opponents ON opponents.match_id = team_results.match_id
				AND opponents.team_key <> team_results.team_key
		),
		stats AS (
			SELECT
				team_key,
				player_ids,
				COUNT(*)::INT AS games_played,
				COUNT(*) FILTER (WHERE is_winner)::INT AS wins,
				COUNT(*) FILTER (WHERE NOT is_winner)::INT AS losses,
				ROUND((COUNT(*) FILTER (WHERE is_winner)::NUMERIC * 100 / COUNT(*)), 2)::FLOAT8 AS win_rate,
				COALESCE(SUM(final_score), 0)::INT AS points_for,
				COALESCE(SUM(opponent_score), 0)::INT AS points_against,
				(COALESCE(SUM(final_score), 0) - COALESCE(SUM(opponent_score), 0))::INT AS point_diff,
				MAX(played_at) AS last_played_at
			FROM team_matchups
			GROUP BY team_key, player_ids
		),
		ranked AS (
			SELECT
				ROW_NUMBER() OVER (
					ORDER BY wins DESC,
						win_rate DESC,
						point_diff DESC,
						points_for DESC,
						games_played DESC,
						team_key ASC
				)::INT AS position,
				team_key,
				player_ids,
				games_played,
				wins,
				losses,
				win_rate,
				points_for,
				points_against,
				point_diff,
				last_played_at
			FROM stats
		)
		SELECT
			position,
			team_key,
			(
				SELECT JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', players.id,
						'name', players.name,
						'avatar_key', players.avatar_key
					)
					ORDER BY team_players.ordinality
				)
				FROM UNNEST(player_ids) WITH ORDINALITY AS team_players(player_id, ordinality)
				INNER JOIN players ON players.id = team_players.player_id
			) AS players,
			games_played,
			wins,
			losses,
			win_rate,
			points_for,
			points_against,
			point_diff,
			last_played_at
		FROM ranked
		ORDER BY position ASC
		LIMIT $2
	`, teamSize, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []TeamEntry
	for rows.Next() {
		var entry TeamEntry
		var playersJSON []byte

		if err := rows.Scan(
			&entry.Position,
			&entry.TeamKey,
			&playersJSON,
			&entry.GamesPlayed,
			&entry.Wins,
			&entry.Losses,
			&entry.WinRate,
			&entry.PointsFor,
			&entry.PointsAgainst,
			&entry.PointDiff,
			&entry.LastPlayedAt,
		); err != nil {
			return nil, err
		}

		if len(playersJSON) > 0 {
			if err := json.Unmarshal(playersJSON, &entry.Players); err != nil {
				return nil, err
			}
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if entries == nil {
		return []TeamEntry{}, nil
	}

	return entries, nil
}

func intPtr(value int) *int {
	return &value
}

func floatPtr(value float64) *float64 {
	return &value
}
