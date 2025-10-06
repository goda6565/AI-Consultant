-- name: CreateHearingMap :exec
INSERT INTO hearing_maps (id, hearing_id, problem_id, content) VALUES ($1, $2, $3, $4);

-- name: DeleteHearingMapByHearingID :execrows
DELETE FROM hearing_maps WHERE hearing_id = $1;

-- name: GetHearingMapByHearingID :one
SELECT * FROM hearing_maps WHERE hearing_id = $1;
