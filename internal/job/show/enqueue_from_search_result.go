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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, '{}'::jsonb)
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
	).Scan(&result.InternalShowID); err != nil {
		return EnqueueFromSearchResultResult{}, err
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
