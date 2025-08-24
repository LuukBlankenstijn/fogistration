-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS wallpaper_configs (
  contest_id  VARCHAR PRIMARY KEY
             REFERENCES contests(external_id) ON DELETE CASCADE,
  filename    TEXT,
  config      JSONB,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- updated_at trigger
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END; $$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_wallpaper_configs_updated_at ON wallpaper_configs;
CREATE TRIGGER trg_wallpaper_configs_updated_at
BEFORE UPDATE ON wallpaper_configs
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

---- create above / drop below ----

DROP TRIGGER IF EXISTS trg_wallpaper_configs_updated_at ON wallpaper_configs;
DROP FUNCTION IF EXISTS set_updated_at();
DROP TABLE IF EXISTS wallpaper_configs;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
