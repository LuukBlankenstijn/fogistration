-- up
CREATE TABLE users (
  id             BIGSERIAL PRIMARY KEY,
  username       TEXT NOT NULL,
  email          TEXT NOT NULL,
  role           TEXT NOT NULL DEFAULT 'user',
  external_id    TEXT,                         -- Keycloak sub (NULL for local)
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_login_at  TIMESTAMPTZ,
  CONSTRAINT users_updated_at_chk CHECK (updated_at IS NOT NULL)
);

-- case-insensitive unique username
CREATE UNIQUE INDEX ux_users_username_ci ON users (lower(username));

-- external_id unique when present
CREATE UNIQUE INDEX ux_users_external_id ON users (external_id) WHERE external_id IS NOT NULL;

-- secrets for local users
CREATE TABLE auth_secrets (
  user_id       BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  password_hash TEXT NOT NULL,
  salt          TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- touch helper
CREATE OR REPLACE FUNCTION tg_touch_updated_at()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END $$;

CREATE TRIGGER users_touch
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION tg_touch_updated_at();

CREATE TRIGGER auth_secrets_touch
BEFORE UPDATE ON auth_secrets
FOR EACH ROW EXECUTE FUNCTION tg_touch_updated_at();

---- create above / drop below ----

-- down
DROP TRIGGER IF EXISTS auth_secrets_touch ON auth_secrets;
DROP TRIGGER IF EXISTS users_touch ON users;
DROP FUNCTION IF EXISTS tg_touch_updated_at();

DROP TABLE IF EXISTS auth_secrets;
DROP INDEX IF EXISTS ux_users_external_id;
DROP INDEX IF EXISTS ux_users_username_ci;
DROP TABLE IF EXISTS users;

