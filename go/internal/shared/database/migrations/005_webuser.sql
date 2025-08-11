-- Write your migrate up statements here
CREATE TABLE app_user (
  id           BIGSERIAL PRIMARY KEY,
  username     TEXT NOT NULL,
  email        TEXT NOT NULL,
  role         TEXT NOT NULL DEFAULT 'user',
  external_id  TEXT,                        -- Keycloak sub (NULL for local)

  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_login_at TIMESTAMPTZ,

  -- ensure updated_at touches
  CONSTRAINT app_user_updated_at_chk CHECK (updated_at IS NOT NULL)
);

-- case-insensitive unique username
CREATE UNIQUE INDEX ux_user_username_ci ON app_user (lower(username));

-- external_id unique when present
CREATE UNIQUE INDEX ux_user_external_id ON app_user (external_id) WHERE external_id IS NOT NULL;

-- secrets for local users
CREATE TABLE auth_secret (
  user_id       BIGINT PRIMARY KEY REFERENCES app_user(id) ON DELETE CASCADE,
  password_hash TEXT NOT NULL,
  salt          TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- touch helper
CREATE OR REPLACE FUNCTION tg_touch_updated_at()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN NEW.updated_at := now(); RETURN NEW; END $$;

CREATE TRIGGER app_user_touch
BEFORE UPDATE ON app_user
FOR EACH ROW EXECUTE FUNCTION tg_touch_updated_at();

CREATE TRIGGER auth_secret_touch
BEFORE UPDATE ON auth_secret
FOR EACH ROW EXECUTE FUNCTION tg_touch_updated_at();

---- create above / drop below ----
DROP TRIGGER IF EXISTS auth_secret_touch ON auth_secret;
DROP TRIGGER IF EXISTS app_user_touch ON app_user;
DROP FUNCTION IF EXISTS tg_touch_updated_at();
DROP TABLE IF EXISTS auth_secret;
DROP INDEX IF EXISTS ux_user_external_id;
DROP INDEX IF EXISTS ux_user_username_ci;
DROP TABLE IF EXISTS app_user;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
