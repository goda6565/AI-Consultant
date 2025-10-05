-- name: GetJobConfigByProblemID :one
SELECT * FROM job_configs WHERE problem_id = $1;

-- name: CreateJobConfig :exec
INSERT INTO job_configs (id, problem_id, enable_internal_search) VALUES ($1, $2, $3);

-- name: UpdateJobConfig :exec
UPDATE job_configs SET enable_internal_search = $2 WHERE id = $1;

-- name: DeleteJobConfigByProblemID :execrows
DELETE FROM job_configs WHERE problem_id = $1;
