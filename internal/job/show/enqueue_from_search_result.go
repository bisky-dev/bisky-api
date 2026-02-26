package show

import (
	"context"

	"github.com/jackc/pgx/v5"
)

const insertShowSQL = `
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, jsonb_build_object('externalId', $13::text))
RETURNING internal_show_id
`

const insertShowJobSQL = `
INSERT INTO job_shows (
  show_id,
  status
)
VALUES ($1::uuid, 'pending')
RETURNING internal_job_show_id, status, retry_count
`

const insertEpisodeSQL = `
INSERT INTO episodes (
  show_id,
  season_number,
  episode_number,
  title,
  air_date,
  runtime_minutes,
  external_ids
)
VALUES ($1::uuid, $2, $3, $4, $5, $6, jsonb_build_object('externalId', $7::text))
ON CONFLICT (show_id, season_number, episode_number) DO NOTHING
`

const findPendingJobByExternalIDSQL = `
SELECT
  s.internal_show_id,
  j.internal_job_show_id,
  j.status,
  j.retry_count
FROM job_shows j
JOIN shows s ON s.internal_show_id = j.show_id
WHERE j.status = 'pending'
  AND s.external_ids->>'externalId' = $1
LIMIT 1
`

func (s *Service) EnqueueFromSearchResult(ctx context.Context, params EnqueueFromSearchResultParams) (EnqueueFromSearchResultResult, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return EnqueueFromSearchResultResult{}, err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	var result EnqueueFromSearchResultResult
	err = tx.QueryRow(ctx, findPendingJobByExternalIDSQL, params.ExternalID).Scan(
		&result.InternalShowID,
		&result.InternalJobShowID,
		&result.Status,
		&result.RetryCount,
	)
	if err == nil {
		if err := tx.Commit(ctx); err != nil {
			return EnqueueFromSearchResultResult{}, err
		}
		committed = true
		return result, nil
	}
	if err != pgx.ErrNoRows {
		return EnqueueFromSearchResultResult{}, err
	}

	if err := tx.QueryRow(ctx, insertShowSQL,
		params.TitlePreferred,
		params.TitleOriginal,
		params.AltTitles,
		params.Type,
		params.Status,
		params.Synopsis,
		params.StartDate,
		params.EndDate,
		params.PosterURL,
		params.BannerURL,
		params.SeasonCount,
		params.EpisodeCount,
		params.ExternalID,
	).Scan(&result.InternalShowID); err != nil {
		return EnqueueFromSearchResultResult{}, err
	}

	for _, episode := range params.Episodes {
		if _, err := tx.Exec(ctx, insertEpisodeSQL,
			result.InternalShowID,
			episode.SeasonNumber,
			episode.EpisodeNumber,
			episode.Title,
			episode.AirDate,
			episode.RuntimeMinutes,
			episode.ExternalID,
		); err != nil {
			return EnqueueFromSearchResultResult{}, err
		}
	}

	if err := tx.QueryRow(ctx, insertShowJobSQL, result.InternalShowID).Scan(&result.InternalJobShowID, &result.Status, &result.RetryCount); err != nil {
		return EnqueueFromSearchResultResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return EnqueueFromSearchResultResult{}, err
	}
	committed = true

	return result, nil
}
