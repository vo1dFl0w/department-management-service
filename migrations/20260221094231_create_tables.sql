-- +goose Up
-- +goose StatementBegin

CREATE TABLE departments (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id BIGINT NULL REFERENCES departments(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT departments_name_length CHECK (char_length(trim(name)) BETWEEN 1 AND 200)
);

CREATE UNIQUE INDEX idx_departments_parent_name_unique
  ON departments ((COALESCE(parent_id, 0)), lower(trim(name)));

CREATE TABLE employees (
    id BIGSERIAL PRIMARY KEY,
    department_id BIGINT NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    full_name TEXT NOT NULL,
    position TEXT NOT NULL,
    hired_at DATE NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT employees_fullname_length CHECK (char_length(trim(full_name)) BETWEEN 1 AND 200),
    CONSTRAINT employees_position_length CHECK (char_length(trim(position)) BETWEEN 1 AND 200)
);

CREATE INDEX idx_employees_department_created_at ON employees (department_id, created_at);
CREATE INDEX idx_departments_parent_id ON departments (parent_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_employees_department_created_at;
DROP INDEX IF EXISTS idx_departments_parent_id;
DROP INDEX IF EXISTS idx_departments_parent_name_unique;

DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;

-- +goose StatementEnd