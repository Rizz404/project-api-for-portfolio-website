-- +goose up
-- * Buat function untuk trigger updated_at saat update
CREATE
OR REPLACE FUNCTION trigger_set_timestamp () RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- +goose down
DROP FUNCTION IF EXISTS trigger_set_timestamp ();
