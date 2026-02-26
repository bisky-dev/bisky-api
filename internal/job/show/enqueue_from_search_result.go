package show

import (
	"context"

	"github.com/jackc/pgx/v5"
)

const insertSearchResultSQL = `
INSERT INTO search_results (
  provider,
  external_id,
  title_preferred,
  title_original,
  type,
  score,
  description,
  banner_url
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING internal_search_result_id
`

const insertShowJobSQL = `
INSERT INTO show_jobs (
  search_result_id,
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
	if err := tx.QueryRow(ctx, insertSearchResultSQL,
		params.Provider,
		params.ExternalID,
		params.TitlePreferred,
		params.TitleOriginal,
		params.Type,
		params.Score,
		params.Description,
		params.BannerURL,
	).Scan(&result.InternalSearchResultID); err != nil {
		return EnqueueFromSearchResultResult{}, err
	}

	if err := tx.QueryRow(ctx, insertShowJobSQL, result.InternalSearchResultID).Scan(&result.InternalJobShowID, &result.Status, &result.RetryCount); err != nil {
		return EnqueueFromSearchResultResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return EnqueueFromSearchResultResult{}, err
	}
	committed = true

	return result, nil
}
