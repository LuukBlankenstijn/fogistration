-- Write your migrate up statements here
--
---- Contests table
CREATE TABLE IF NOT EXISTS contests (
    id INTEGER PRIMARY KEY,
    external_id VARCHAR NOT NULL,
    formal_name VARCHAR NOT NULL,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    hash VARCHAR NOT NULL
);

-- Teams table  
CREATE TABLE IF NOT EXISTS teams (
    id INTEGER PRIMARY KEY,
    external_id VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    display_name VARCHAR,
    ip VARCHAR,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    hash VARCHAR NOT NULL
);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    team_id INTEGER REFERENCES teams (id),
    ip VARCHAR,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    hash VARCHAR NOT NULL
);

-- Contest-Team relationship (many-to-many)
CREATE TABLE IF NOT EXISTS contest_teams (
    contest_id INTEGER REFERENCES contests (id) ON DELETE CASCADE,
    team_id INTEGER REFERENCES teams (id) ON DELETE CASCADE,
    PRIMARY KEY (contest_id, team_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_team_id ON users (team_id);
CREATE INDEX IF NOT EXISTS idx_contest_teams_contest_id ON contest_teams (contest_id);
CREATE INDEX IF NOT EXISTS idx_contest_teams_team_id ON contest_teams (team_id);

---- create above / drop below ----
--
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

-- Drop indexes
DROP INDEX IF EXISTS idx_contest_teams_team_id;
DROP INDEX IF EXISTS idx_contest_teams_contest_id;
DROP INDEX IF EXISTS idx_users_team_id;

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS contest_teams;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS contests;

