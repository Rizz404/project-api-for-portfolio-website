-- +goose up
CREATE TABLE tech_stacks (
  id_project UUID NOT NULL,
  id_tech UUID NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  PRIMARY KEY (id_project, id_tech),
  -- * Jika project dihapus, relasi ini ikut hilang
  CONSTRAINT fk_tech_stacks_project FOREIGN KEY (id_project) REFERENCES projects (id) ON DELETE CASCADE,
  -- * Jika tech dihapus, relasi ini ikut hilang
  CONSTRAINT fk_tech_stacks_tech FOREIGN KEY (id_tech) REFERENCES techs (id) ON DELETE CASCADE
);

CREATE TRIGGER set_timestamp BEFORE
UPDATE ON tech_stacks FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp ();

CREATE INDEX idx_tech_stacks_tech_id ON tech_stacks (id_tech);

-- +goose down
DROP TRIGGER IF EXISTS set_timestamp ON tech_stacks;

DROP TABLE IF EXISTS tech_stacks;
