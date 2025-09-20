DROP INDEX IF EXISTS idx_problem_fields_problem_id;

DROP TABLE IF EXISTS problem_fields;

ALTER TABLE hearing_messages DROP COLUMN problem_field_id;