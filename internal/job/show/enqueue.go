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
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	pending, found, err := findPendingJobByExternalID(ctx, tx, params.ExternalID)
	if err != nil {
		return EnqueueFromSearchResultResult{}, err
	}
	if found {
		if err := tx.Commit(ctx); err != nil {
			return EnqueueFromSearchResultResult{}, err
		}
		return pending, nil
	}

	internalShowID, err := insertShow(ctx, tx, params)
	if err != nil {
		return EnqueueFromSearchResultResult{}, err
	}

	if err := insertEpisodes(ctx, tx, internalShowID, params.Episodes); err != nil {
		return EnqueueFromSearchResultResult{}, err
	}

	result, err := insertShowJob(ctx, tx, internalShowID)
	if err != nil {
		return EnqueueFromSearchResultResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return EnqueueFromSearchResultResult{}, err
	}
	return result, nil
}

func findPendingJobByExternalID(ctx context.Context, tx pgx.Tx, externalID string) (EnqueueFromSearchResultResult, bool, error) {
	var result EnqueueFromSearchResultResult
	err := tx.QueryRow(ctx, findPendingJobByExternalIDSQL, externalID).Scan(
		&result.InternalShowID,
		&result.InternalJobShowID,
		&result.Status,
		&result.RetryCount,
	)
	if err == nil {
		return result, true, nil
	}
	if err == pgx.ErrNoRows {
		return EnqueueFromSearchResultResult{}, false, nil
	}
	return EnqueueFromSearchResultResult{}, false, err
}

func insertShow(ctx context.Context, tx pgx.Tx, params EnqueueFromSearchResultParams) (string, error) {
	var internalShowID string
	err := tx.QueryRow(ctx, insertShowSQL,
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
	).Scan(&internalShowID)
	if err != nil {
		return "", err
	}
	return internalShowID, nil
}

func insertEpisodes(ctx context.Context, tx pgx.Tx, internalShowID string, episodes []EpisodeInput) error {
	for _, episode := range episodes {
		if _, err := tx.Exec(ctx, insertEpisodeSQL,
			internalShowID,
			episode.SeasonNumber,
			episode.EpisodeNumber,
			episode.Title,
			episode.AirDate,
			episode.RuntimeMinutes,
			episode.ExternalID,
		); err != nil {
			return err
		}
	}
	return nil
}

func insertShowJob(ctx context.Context, tx pgx.Tx, internalShowID string) (EnqueueFromSearchResultResult, error) {
	result := EnqueueFromSearchResultResult{InternalShowID: internalShowID}
	err := tx.QueryRow(ctx, insertShowJobSQL, internalShowID).Scan(
		&result.InternalJobShowID,
		&result.Status,
		&result.RetryCount,
	)
	if err != nil {
		return EnqueueFromSearchResultResult{}, err
	}
	return result, nil
}
