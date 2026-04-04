CREATE TABLE player_stats (
  player_id    UUID REFERENCES players(id) ON DELETE CASCADE,
  mode         TEXT NOT NULL,
  games_played INT DEFAULT 0,
  wins         INT DEFAULT 0,
  losses       INT DEFAULT 0,
  updated_at   TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (player_id, mode)
);

CREATE INDEX idx_player_stats_mode ON player_stats(mode);