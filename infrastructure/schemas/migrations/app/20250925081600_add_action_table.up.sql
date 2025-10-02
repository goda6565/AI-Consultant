CREATE TABLE actions (
    id UUID PRIMARY KEY,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    action_type TEXT NOT NULL,
    input TEXT NOT NULL,
    output TEXT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_actions_problem_id ON actions(problem_id);
