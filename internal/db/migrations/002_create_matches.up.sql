CREATE TABLE matches (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  mode        TEXT NOT NULL,
  status      TEXT NOT NULL DEFAULT 'ongoing',
  room_code   TEXT NOT NULL UNIQUE,
  max_score   INT,
  created_at  TIMESTAMPTZ DEFAULT NOW(),
  finished_at TIMESTAMPTZ
);

CREATE INDEX idx_matches_status ON matches(status);
CREATE INDEX idx_matches_room_code ON matches(room_code);