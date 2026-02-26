ALTER INDEX uq_show_jobs_pending RENAME TO uq_jobs_show_pending;
ALTER INDEX idx_show_jobs_status_created_at RENAME TO idx_jobs_show_status_created_at;
ALTER INDEX idx_search_results_provider_external_id RENAME TO idx_search_result_provider_external_id;

ALTER TABLE show_jobs RENAME TO jobs_show;
ALTER TABLE search_results RENAME TO search_result;
