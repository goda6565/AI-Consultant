-- name: CreateProblem :exec
INSERT INTO problems (id, title, description, status) VALUES ($1, $2, $3, $4);

-- name: DeleteProblem :execrows
DELETE FROM problems WHERE id = $1;

-- name: GetProblemById :one
SELECT * FROM problems WHERE id = $1;

-- name: GetAllProblems :many
SELECT * FROM problems ORDER BY created_at DESC;

-- name: UpdateProblemStatus :execrows
UPDATE problems SET status = $2 WHERE id = $1;