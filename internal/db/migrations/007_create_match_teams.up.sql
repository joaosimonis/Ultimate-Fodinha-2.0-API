CREATE TABLE match_teams (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  match_id      UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
  name          TEXT NOT NULL,
  display_order INT NOT NULL DEFAULT 0,
  UNIQUE (match_id, name)
);

CREATE INDEX idx_match_teams_match_id ON match_teams(match_id);

ALTER TABLE match_players
  ADD COLUMN team_id UUID REFERENCES match_teams(id) ON DELETE SET NULL;

WITH distinct_teams AS (
  SELECT
    match_id,
    team_name AS name,
    ROW_NUMBER() OVER (PARTITION BY match_id ORDER BY team_name)::INT AS display_order
  FROM (
    SELECT DISTINCT match_id, team_name
    FROM match_players
    WHERE NULLIF(TRIM(team_name), '') IS NOT NULL
  ) teams
)
INSERT INTO match_teams (match_id, name, display_order)
SELECT match_id, name, display_order
FROM distinct_teams
ON CONFLICT (match_id, name) DO NOTHING;

UPDATE match_players mp
SET team_id = mt.id
FROM match_teams mt
WHERE mt.match_id = mp.match_id
  AND mt.name = mp.team_name
  AND NULLIF(TRIM(mp.team_name), '') IS NOT NULL;

CREATE INDEX idx_match_players_team_id ON match_players(team_id);
