-- name: CreateHearing :exec
INSERT INTO hearings (id, problem_id) VALUES ($1, $2);

-- name: DeleteHearingByProblemID :execrows
DELETE FROM hearings WHERE problem_id = $1;

-- name: GetHearingById :one
SELECT * FROM hearings WHERE id = $1;

-- name: GetHearingByProblemId :one
SELECT * FROM hearings WHERE problem_id = $1;

-- name: GetAllHearingsByProblemId :many
SELECT * FROM hearings WHERE problem_id = $1;