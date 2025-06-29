-- +goose up
CREATE TABLE techs (
  id UUID NOT NULL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  description TEXT,
  logo_url TEXT,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_techs_name ON techs (name);

CREATE TRIGGER set_timestamp BEFORE
UPDATE ON techs FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose down
DROP TRIGGER IF EXISTS set_timestamp ON techs;

DROP INDEX IF EXISTS idx_techs_name;

DROP TABLE IF EXISTS techs;
