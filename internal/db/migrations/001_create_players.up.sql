
CREATE TABLE players (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name       TEXT NOT NULL UNIQUE,
  avatar_key TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);