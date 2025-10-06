CREATE TABLE hearing_maps (
    id UUID PRIMARY KEY,
    hearing_id UUID NOT NULL REFERENCES hearings(id) ON DELETE RESTRICT,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE RESTRICT,
    content TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_hearing_maps_hearing_id ON hearing_maps(hearing_id);
CREATE INDEX idx_hearing_maps_problem_id ON hearing_maps(problem_id);
