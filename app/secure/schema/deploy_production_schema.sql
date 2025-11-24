-- deploy_production_schema.sql
DROP TABLE IF EXISTS game_state;
DROP TABLE IF EXISTS users;

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    fingerprint     TEXT NOT NULL UNIQUE CHECK (trim(fingerprint) != ''),
    ip              TEXT,
    user_agent      TEXT,
    device_type     TEXT,
    os              TEXT,
    creation_date   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS game_state (
    user_id         INTEGER NOT NULL,
    state           TEXT NOT NULL, -- serialized JSON
    date            DATE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, date),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
