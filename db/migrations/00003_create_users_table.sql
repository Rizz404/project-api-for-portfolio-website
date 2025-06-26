-- +goose Up
CREATE TYPE user_role AS ENUM('admin', 'user');

CREATE TABLE
  users (
    id UUID NOT NULL PRIMARY KEY,
    username VARCHAR(30) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    address TEXT,
    full_name TEXT,
    id_language UUID,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
  );

-- * Constraintnya
ALTER TABLE users
ADD CONSTRAINT fk_users_language FOREIGN KEY (id_language) REFERENCES languages (id) ON DELETE SET NULL;

CREATE INDEX idx_users_username ON users (username);

CREATE INDEX idx_users_email ON users (email);

-- * Pake func buat trigger timestamp
CREATE TRIGGER set_timestamp BEFORE
UPDATE ON users FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose Down
DROP TRIGGER IF EXISTS set_timestamp ON users;

DROP INDEX IF EXISTS idx_users_username;

DROP INDEX IF EXISTS idx_users_email;

ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_users_language;

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS user_role;
