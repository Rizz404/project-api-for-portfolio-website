-- +goose up
CREATE TABLE
  projects (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    id_category UUID NOT NULL,
    description TEXT,
    is_deployed BOOLEAN NOT NULL DEFAULT false,
    is_maintained BOOLEAN NOT NULL DEFAULT true,
    live_demo VARCHAR(255),
    source_code VARCHAR(255),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
  );

-- * Constraintnya
ALTER TABLE projects
ADD CONSTRAINT fk_projects_category FOREIGN KEY (id_category) REFERENCES categories (id) ON DELETE RESTRICT;

-- * Indexnya
CREATE INDEX idx_projects_name ON projects (name);

CREATE INDEX idx_projects_id_category ON projects (id_category);

CREATE TRIGGER set_timestamp BEFORE
UPDATE ON projects FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose down
DROP TRIGGER IF EXISTS set_timestamp ON projects;

DROP INDEX IF EXISTS idx_projects_name;

DROP TABLE IF EXISTS projects;
