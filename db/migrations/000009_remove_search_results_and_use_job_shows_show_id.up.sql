CREATE TABLE _tmp_search_result_show_map (
  search_result_id UUID PRIMARY KEY,
  show_id UUID NOT NULL UNIQUE
);

INSERT INTO _tmp_search_result_show_map (search_result_id, show_id)
SELECT internal_search_result_id, gen_random_uuid()
FROM search_results;

INSERT INTO shows (
  internal_show_id,
  title_preferred,
  title_original,
  alt_titles,
  type,
  status,
  synopsis,
  start_date,
  end_date,
  poster_url,
  banner_url,
  season_count,
  episode_count,
  external_ids
)
SELECT
  m.show_id,
  sr.title_preferred,
  sr.title_original,
  '{}'::TEXT[],
  CASE
    WHEN sr.type IN ('anime', 'tv', 'movie', 'ova', 'special') THEN sr.type
    ELSE 'tv'
  END,
  'ongoing',
  sr.description,
  NULL,
  NULL,
  NULL,
  sr.banner_url,
  NULL,
  NULL,
  '{}'::jsonb
FROM search_results sr
JOIN _tmp_search_result_show_map m ON m.search_result_id = sr.internal_search_result_id;

ALTER TABLE show_jobs
ADD COLUMN show_id UUID;

UPDATE show_jobs sj
SET show_id = m.show_id
FROM _tmp_search_result_show_map m
WHERE sj.search_result_id = m.search_result_id;

DROP INDEX IF EXISTS uq_show_jobs_pending;

ALTER TABLE show_jobs
DROP CONSTRAINT IF EXISTS show_jobs_search_result_id_fkey;

ALTER TABLE show_jobs
DROP COLUMN search_result_id;

ALTER TABLE show_jobs
ALTER COLUMN show_id SET NOT NULL;

ALTER TABLE show_jobs
ADD CONSTRAINT show_jobs_show_id_fkey
FOREIGN KEY (show_id) REFERENCES shows(internal_show_id) ON DELETE CASCADE;

ALTER TABLE show_jobs RENAME TO job_shows;

ALTER INDEX idx_show_jobs_status_created_at RENAME TO idx_job_shows_status_created_at;

CREATE UNIQUE INDEX uq_job_shows_pending ON job_shows (show_id, status)
WHERE status IN ('pending', 'processing');

DROP TABLE search_results;
DROP TABLE _tmp_search_result_show_map;
