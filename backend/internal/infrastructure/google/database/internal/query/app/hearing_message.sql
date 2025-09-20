-- name: CreateHearingMessage :exec
INSERT INTO hearing_messages (id, hearing_id, problem_field_id, role, message) VALUES ($1, $2, $3, $4, $5);

-- name: DeleteHearingMessageByHearingID :execrows
DELETE FROM hearing_messages WHERE hearing_id = $1;

-- name: GetHearingMessageByHearingID :many
SELECT *
FROM hearing_messages
WHERE hearing_id = $1
ORDER BY created_at ASC, id ASC;

-- name: GetHearingMessageByProblemFieldID :many
SELECT *
FROM hearing_messages
WHERE problem_field_id = $1
ORDER BY created_at ASC, id ASC;
