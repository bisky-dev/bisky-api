CREATE TABLE search_results (
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

CREATE INDEX idx_search_results_provider_external_id ON search_results (provider, external_id);

ALTER TABLE job_shows RENAME TO show_jobs;

ALTER INDEX idx_job_shows_status_created_at RENAME TO idx_show_jobs_status_created_at;

DROP INDEX IF EXISTS uq_job_shows_pending;

ALTER TABLE show_jobs
ADD COLUMN search_result_id UUID;

CREATE TABLE _tmp_show_search_result_map (
  show_id UUID PRIMARY KEY,
  search_result_id UUID NOT NULL UNIQUE
);

INSERT INTO _tmp_show_search_result_map (show_id, search_result_id)
SELECT DISTINCT show_id, gen_random_uuid()
FROM show_jobs;

INSERT INTO search_results (
  internal_search_result_id,
  provider,
  external_id,
  title_preferred,
  title_original,
  type,
  description,
  banner_url
)
SELECT
  m.search_result_id,
  'anidb',
  s.internal_show_id::TEXT,
  s.title_preferred,
  s.title_original,
  s.type,
  s.synopsis,
  s.banner_url
FROM _tmp_show_search_result_map m
JOIN shows s ON s.internal_show_id = m.show_id;

UPDATE show_jobs sj
SET search_result_id = m.search_result_id
FROM _tmp_show_search_result_map m
WHERE sj.show_id = m.show_id;

ALTER TABLE show_jobs
DROP CONSTRAINT IF EXISTS show_jobs_show_id_fkey;

ALTER TABLE show_jobs
DROP COLUMN show_id;

ALTER TABLE show_jobs
ALTER COLUMN search_result_id SET NOT NULL;

ALTER TABLE show_jobs
ADD CONSTRAINT show_jobs_search_result_id_fkey
FOREIGN KEY (search_result_id) REFERENCES search_results(internal_search_result_id) ON DELETE CASCADE;

CREATE UNIQUE INDEX uq_show_jobs_pending ON show_jobs (search_result_id, status)
WHERE status IN ('pending', 'processing');

DROP TABLE _tmp_show_search_result_map;
