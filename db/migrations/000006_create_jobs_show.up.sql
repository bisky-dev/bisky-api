CREATE TABLE jobs_show (
  internal_job_show_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  search_result_id UUID NOT NULL REFERENCES search_result(internal_search_result_id) ON DELETE CASCADE,
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
  error_message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_jobs_show_status_created_at ON jobs_show (status, created_at);
CREATE UNIQUE INDEX uq_jobs_show_pending ON jobs_show (search_result_id, status)
WHERE status IN ('pending', 'processing');
