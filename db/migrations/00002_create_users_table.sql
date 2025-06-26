-- +goose up
CREATE TABLE
  users (
    id UUID NOT NULL PRIMARY KEY,
    username VARCHAR(30) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    address TEXT,
    full_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
  );

CREATE INDEX idx_users_username ON users (username);

CREATE INDEX idx_users_email ON users (email);

-- * Pake func buat trigger timestamp
CREATE TRIGGER set_timestamp BEFORE
UPDATE ON users FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp ();

-- +goose down
DROP TABLE IF EXISTS users;

DROP TRIGGER IF EXISTS set_timestamp ON users;

DROP INDEX IF EXISTS idx_users_username;

DROP INDEX IF EXISTS idx_users_email;

DROP TABLE IF EXISTS users;
