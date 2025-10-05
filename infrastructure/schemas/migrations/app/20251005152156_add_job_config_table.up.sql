CREATE TABLE job_configs (
    id VARCHAR(36) PRIMARY KEY,
    problem_id VARCHAR(36) NOT NULL,
    enable_internal_search BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_job_configs_problem_id ON job_configs (problem_id);
