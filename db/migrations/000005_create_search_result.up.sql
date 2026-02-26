CREATE TABLE search_result (
  internal_search_result_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  provider TEXT NOT NULL CHECK (provider IN ('anidb', 'anilist', 'tvdb')),
  external_id TEXT NOT NULL,
  title_preferred TEXT NOT NULL,
  title_original TEXT,
  type TEXT,
  score DOUBLE PRECISION,
  description TEXT,
  banner_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_search_result_provider_external_id ON search_result (provider, external_id);
