CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE vectors (
    id UUID NOT NULL,
    document_id UUID NOT NULL,
    content TEXT NOT NULL,
    parent_content TEXT NOT NULL,
    embedding vector(1536) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_vectors_embedding ON vectors USING hnsw(embedding vector_cosine_ops) WITH (m = 24, ef_construction = 100);