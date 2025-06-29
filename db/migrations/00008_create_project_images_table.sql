-- +goose up
CREATE TABLE
  project_images (
    id UUID NOT NULL PRIMARY KEY,
    id_project UUID NOT NULL,
    file_name TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
  );

ALTER TABLE project_images
ADD CONSTRAINT fk_project_images_project FOREIGN KEY (id_project) REFERENCES projects (id) ON DELETE RESTRICT;

CREATE INDEX idx_project_images_id_project ON project_images (id_project);

CREATE TRIGGER set_timestamp BEFORE
UPDATE ON project_images FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose down
DROP TRIGGER IF EXISTS set_timestamp ON project_images;

DROP INDEX IF EXISTS idx_project_images_id_project;

DROP TABLE IF EXISTS project_images;
