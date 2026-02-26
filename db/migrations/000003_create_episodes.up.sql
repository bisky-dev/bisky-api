CREATE TABLE episodes (
  internal_episode_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  show_id UUID NOT NULL REFERENCES shows(internal_show_id) ON DELETE CASCADE,
  season_number BIGINT NOT NULL,
  episode_number BIGINT NOT NULL,
  title TEXT NOT NULL,
  air_date TEXT,
  runtime_minutes BIGINT,
  external_ids JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT episodes_season_number_non_negative CHECK (season_number >= 0),
  CONSTRAINT episodes_episode_number_non_negative CHECK (episode_number >= 0),
  CONSTRAINT episodes_runtime_minutes_non_negative CHECK (runtime_minutes IS NULL OR runtime_minutes >= 0),
  CONSTRAINT episodes_show_season_episode_unique UNIQUE (show_id, season_number, episode_number)
);

CREATE INDEX idx_episodes_show_id ON episodes (show_id);
CREATE INDEX idx_episodes_show_id_season_episode ON episodes (show_id, season_number, episode_number);
