-- name: CreateVector :exec
INSERT INTO vectors (id, document_id, content, parent_content, embedding) VALUES ($1, $2, $3, $4, $5);

-- name: SearchVector :many
SELECT id, document_id, content, parent_content, (1 - (embedding <=> $1))::float8 AS similarity FROM vectors ORDER BY similarity DESC LIMIT $2;

-- name: DeleteVector :execrows
DELETE FROM vectors WHERE document_id = $1;