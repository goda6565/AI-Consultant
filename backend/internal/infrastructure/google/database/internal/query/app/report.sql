-- name: CreateReport :exec
INSERT INTO reports (id, problem_id, content) VALUES ($1, $2, $3);

-- name: GetReportByProblemID :one
SELECT * FROM reports WHERE problem_id = $1;

-- name: DeleteReportsByProblemID :execrows
DELETE FROM reports WHERE problem_id = $1;
