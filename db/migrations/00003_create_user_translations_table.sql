-- +goose up
CREATE TABLE
  user_translations (
    id UUID NOT NULL PRIMARY KEY,
    id_user UUID NOT NULL,
    bio TEXT,
    about_me TEXT,
    additional_skills TEXT[] DEFAULT '{}',
    languages TEXT[] DEFAULT '{}',
    quote TEXT,
    lang_code VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
  );

-- * Constraintnya
ALTER TABLE user_translations
ADD CONSTRAINT fk_user_translations_user FOREIGN KEY (id_user) REFERENCES users (id) ON DELETE RESTRICT;

-- * Index
-- * Index GIN untuk pencarian efisien di dalam array
CREATE INDEX idx_user_translations_skills ON user_translations USING GIN (additional_skills);

CREATE INDEX idx_user_translations_languages ON user_translations USING GIN (languages);

CREATE TRIGGER set_timestamp BEFORE
UPDATE ON user_translations FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose down
DROP TRIGGER IF EXISTS set_timestamp ON user_translations;

DROP TABLE IF EXISTS user_translations;
