CREATE TABLE documents (
    id UUID PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    document_extension VARCHAR(50) NOT NULL,
    bucket_name VARCHAR(50) NOT NULL,
    object_name VARCHAR(50) NOT NULL,
    document_status VARCHAR(50) NOT NULL,
    sync_step VARCHAR(50) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_documents_title UNIQUE (title),
    CONSTRAINT uq_documents_storage_path UNIQUE (bucket_name, object_name)
);

CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_documents_set_updated_at
BEFORE UPDATE ON documents
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

