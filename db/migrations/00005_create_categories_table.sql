-- +goose up
CREATE TABLE categories (
  id UUID NOT NULL PRIMARY KEY,
  name VARCHAR(30) NOT NULL UNIQUE,
  description TEXT,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_categories_name ON categories (name);

-- * Pake func buat trigger timestamp
CREATE TRIGGER set_timestamp BEFORE
UPDATE ON categories FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose down
DROP TABLE IF EXISTS categories;

DROP TRIGGER IF EXISTS set_timestamp ON categories;

DROP INDEX IF EXISTS idx_categories_name;

DROP TABLE IF EXISTS categories;
