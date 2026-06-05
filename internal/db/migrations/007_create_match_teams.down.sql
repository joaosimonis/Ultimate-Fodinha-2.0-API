DROP INDEX IF EXISTS idx_match_players_team_id;

ALTER TABLE match_players
  DROP COLUMN IF EXISTS team_id;

DROP INDEX IF EXISTS idx_match_teams_match_id;
DROP TABLE IF EXISTS match_teams;
