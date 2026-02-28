CREATE TABLE hook_settings (
  event_name TEXT PRIMARY KEY,
  target_url TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
