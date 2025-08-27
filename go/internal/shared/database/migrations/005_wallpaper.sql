-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS wallpapers (
  id  INTEGER PRIMARY KEY
             REFERENCES contests(id) ON DELETE CASCADE,
  filename    TEXT,
  layout      JSONB,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- updated_at trigger
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END; $$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_wallpapers_updated_at ON wallpapers;
CREATE TRIGGER trg_wallpapers_updated_at
BEFORE UPDATE ON wallpapers
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

---- create above / drop below ----

DROP TRIGGER IF EXISTS trg_wallpapers_updated_at ON wallpapers;
DROP FUNCTION IF EXISTS set_updated_at();
DROP TABLE IF EXISTS wallpapers;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
