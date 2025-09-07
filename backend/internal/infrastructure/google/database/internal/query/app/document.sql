-- name: GetDocument :one
SELECT * FROM documents WHERE id = $1;

-- name: GetAllDocuments :many
SELECT * FROM documents ORDER BY created_at DESC;

-- name: GetDocumentByTitle :one
SELECT * FROM documents WHERE title = $1;

-- name: CreateDocument :exec
INSERT INTO documents (id, title, document_extension, bucket_name, object_name, document_status, sync_step) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateDocument :execrows
UPDATE documents SET title = $2, document_extension = $3, bucket_name = $4, object_name = $5, document_status = $6, sync_step = $7 WHERE id = $1;

-- name: DeleteDocument :execrows
DELETE FROM documents WHERE id = $1;