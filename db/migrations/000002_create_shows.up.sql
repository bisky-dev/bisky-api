CREATE TABLE shows (
  internal_show_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title_preferred TEXT NOT NULL,
  title_original TEXT,
  alt_titles TEXT[] NOT NULL DEFAULT '{}',
  type TEXT NOT NULL CHECK (type IN ('anime', 'tv', 'movie', 'ova', 'special')),
  status TEXT NOT NULL CHECK (status IN ('ongoing', 'finished')),
  synopsis TEXT,
  start_date TEXT,
  end_date TEXT,
  poster_url TEXT,
  banner_url TEXT,
  season_count BIGINT,
  episode_count BIGINT,
  external_ids JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shows_title_preferred ON shows (title_preferred);
CREATE INDEX idx_shows_type_status ON shows (type, status);
