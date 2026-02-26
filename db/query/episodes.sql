-- name: CreateEpisode :one
INSERT INTO episodes (
  show_id,
  season_number,
  episode_number,
  title,
  air_date,
  runtime_minutes,
  external_ids
)
VALUES ($1::uuid, $2, $3, $4, $5, $6, $7)
RETURNING
  internal_episode_id,
  show_id,
  season_number,
  episode_number,
  title,
  air_date,
  runtime_minutes,
  external_ids,
  created_at,
  updated_at;

-- name: ListEpisodes :many
SELECT
  internal_episode_id,
  show_id,
  season_number,
  episode_number,
  title,
  air_date,
  runtime_minutes,
  external_ids,
  created_at,
  updated_at
FROM episodes
ORDER BY created_at DESC;

-- name: ListEpisodesByShowID :many
SELECT
  internal_episode_id,
  show_id,
  season_number,
  episode_number,
  title,
  air_date,
  runtime_minutes,
  external_ids,
  created_at,
  updated_at
FROM episodes
WHERE show_id = $1::uuid
ORDER BY season_number ASC, episode_number ASC;

-- name: GetEpisodeByID :one
SELECT
  internal_episode_id,
  show_id,
  season_number,
  episode_number,
  title,
  air_date,
  runtime_minutes,
  external_ids,
  created_at,
  updated_at
FROM episodes
WHERE internal_episode_id = $1::uuid
LIMIT 1;

-- name: UpdateEpisode :one
UPDATE episodes
SET
  show_id = $2::uuid,
  season_number = $3,
  episode_number = $4,
  title = $5,
  air_date = $6,
  runtime_minutes = $7,
  external_ids = $8,
  updated_at = NOW()
WHERE internal_episode_id = $1::uuid
RETURNING
  internal_episode_id,
  show_id,
  season_number,
  episode_number,
  title,
  air_date,
  runtime_minutes,
  external_ids,
  created_at,
  updated_at;

-- name: DeleteEpisode :one
DELETE FROM episodes
WHERE internal_episode_id = $1::uuid
RETURNING internal_episode_id;
