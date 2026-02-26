-- name: CreateShow :one
INSERT INTO shows (
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING
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
  external_ids,
  created_at,
  updated_at;

-- name: ListShows :many
SELECT
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
  external_ids,
  created_at,
  updated_at
FROM shows
ORDER BY created_at DESC;

-- name: GetShowByID :one
SELECT
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
  external_ids,
  created_at,
  updated_at
FROM shows
WHERE internal_show_id = $1::uuid
LIMIT 1;

-- name: UpdateShow :one
UPDATE shows
SET
  title_preferred = $2,
  title_original = $3,
  alt_titles = $4,
  type = $5,
  status = $6,
  synopsis = $7,
  start_date = $8,
  end_date = $9,
  poster_url = $10,
  banner_url = $11,
  season_count = $12,
  episode_count = $13,
  external_ids = $14,
  updated_at = NOW()
WHERE internal_show_id = $1::uuid
RETURNING
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
  external_ids,
  created_at,
  updated_at;

-- name: DeleteShow :one
DELETE FROM shows
WHERE internal_show_id = $1::uuid
RETURNING internal_show_id;
