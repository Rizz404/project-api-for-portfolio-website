-- +goose up
CREATE TABLE
  languages (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(30) NOT NULL UNIQUE,
    lang_code VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
  );

CREATE INDEX idx_languages_name ON languages (name);

CREATE INDEX idx_languages_lang_code ON languages (lang_code);

-- * Pake func buat trigger timestamp
CREATE TRIGGER set_timestamp BEFORE
UPDATE ON languages FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose down
DROP TABLE IF EXISTS languages;

DROP TRIGGER IF EXISTS set_timestamp ON languages;

DROP INDEX IF EXISTS idx_languages_name;

DROP INDEX IF EXISTS idx_languages_lang_code;

DROP TABLE IF EXISTS languages;
