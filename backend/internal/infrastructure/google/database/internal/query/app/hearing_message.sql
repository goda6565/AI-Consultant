-- name: CreateHearingMessage :exec
INSERT INTO hearing_messages (id, hearing_id, role, message) VALUES ($1, $2, $3, $4);

-- name: DeleteHearingMessageByHearingID :execrows
DELETE FROM hearing_messages WHERE hearing_id = $1;

-- name: GetHearingMessageByHearingID :many
SELECT * FROM hearing_messages WHERE hearing_id = $1;