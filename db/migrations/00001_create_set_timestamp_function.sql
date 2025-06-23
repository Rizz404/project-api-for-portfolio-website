-- +goose up
-- * Anotasi untuk memberitahu goose agar tidak memecah statement ini
-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION trigger_set_timestamp () RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose down
DROP FUNCTION IF EXISTS trigger_set_timestamp ();
