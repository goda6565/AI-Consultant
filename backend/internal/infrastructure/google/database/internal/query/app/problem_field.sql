-- name: CreateProblemField :exec
INSERT INTO problem_fields (id, problem_id, field, answered) VALUES ($1, $2, $3, $4);

-- name: UpdateAnswered :execrows
UPDATE problem_fields SET answered = $2 WHERE id = $1;

-- name: DeleteProblemField :execrows
DELETE FROM problem_fields WHERE id = $1;

-- name: FindByProblemID :many
SELECT * FROM problem_fields WHERE problem_id = $1;

