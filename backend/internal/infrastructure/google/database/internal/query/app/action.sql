-- name: CreateAction :exec
INSERT INTO actions (id, problem_id, action_type, input, output) VALUES ($1, $2, $3, $4, $5);

-- name: DeleteActionsByProblemID :execrows
DELETE FROM actions WHERE problem_id = $1;

-- name: GetActionsByProblemID :many
SELECT * FROM actions WHERE problem_id = $1 ORDER BY created_at ASC;
