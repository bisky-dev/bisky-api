ALTER TABLE search_result RENAME TO search_results;
ALTER TABLE jobs_show RENAME TO show_jobs;

ALTER INDEX idx_search_result_provider_external_id RENAME TO idx_search_results_provider_external_id;
ALTER INDEX idx_jobs_show_status_created_at RENAME TO idx_show_jobs_status_created_at;
ALTER INDEX uq_jobs_show_pending RENAME TO uq_show_jobs_pending;
