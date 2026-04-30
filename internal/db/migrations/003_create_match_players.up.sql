CREATE TABLE match_players (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  match_id    UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
  player_id   UUID NOT NULL REFERENCES players(id),
  team_name   TEXT,
  final_score INT,
  remaining_lives INT,
  placement   INT,
  is_winner   BOOLEAN DEFAULT FALSE,
  UNIQUE (match_id, player_id)
);

CREATE INDEX idx_match_players_match_id ON match_players(match_id);
