CREATE TABLE players (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name       TEXT NOT NULL UNIQUE,
  avatar_key TEXT NOT NULL,
  title      TEXT,
  lore       TEXT,
  ultimate   JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);