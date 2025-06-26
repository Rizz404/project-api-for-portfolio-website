-- +goose up
CREATE TABLE
  project_translations (
    id UUID NOT NULL PRIMARY KEY,
    id_project UUID NOT NULL,
    id_language UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
  );

-- * Constraintnya
ALTER TABLE project_translations
ADD CONSTRAINT fk_project_translations_project FOREIGN KEY (id_project) REFERENCES projects (id) ON DELETE RESTRICT;

ALTER TABLE project_translations
ADD CONSTRAINT fk_project_translations_language FOREIGN KEY (id_language) REFERENCES languages (id) ON DELETE RESTRICT;

-- * Indexnya
CREATE INDEX idx_project_translations_name ON project_translations (name);

CREATE TRIGGER set_timestamp BEFORE
UPDATE ON project_translations FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose down
DROP TRIGGER IF EXISTS set_timestamp ON project_translations;

DROP TABLE IF EXISTS project_translations;
