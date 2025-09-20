CREATE TABLE problem_fields (
    id UUID PRIMARY KEY,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE RESTRICT,
    field TEXT NOT NULL,
    answered BOOLEAN NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE hearing_messages ADD problem_field_id UUID REFERENCES problem_fields(id) ON DELETE RESTRICT;

CREATE INDEX idx_problem_fields_problem_id ON problem_fields(problem_id);