-- deploy_production_schema.sql
-- DROP TABLE IF EXISTS game_state;
-- DROP TABLE IF EXISTS users;

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    cognito_sub         TEXT,
    cognito_username    TEXT,
    fingerprint         TEXT NOT NULL UNIQUE CHECK (trim(fingerprint) != ''),
    ip                  TEXT,
    user_agent          TEXT,
    device_type         TEXT,
    os                  TEXT,
    creation_date       DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS game_state (
    user_id         INTEGER NOT NULL,
    state           TEXT NOT NULL, -- serialized JSON
    date            DATE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, date),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS login_streaks (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL,
    login_date DATE NOT NULL DEFAULT (DATE('now', 'localtime')),
    streak     INTEGER NOT NULL DEFAULT 1,
    UNIQUE (user_id, login_date),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- CREATE INDEX IF NOT EXISTS idx_login_streaks_user_date
-- ON login_streaks (user_id, login_date);

CREATE TABLE IF NOT EXISTS word_chain_game_state (
    user_id      INTEGER NOT NULL,
    puzzle_date  DATE NOT NULL,
    guesses      TEXT NOT NULL, -- JSON array of {"word": string, "result": boolean}
    PRIMARY KEY (user_id, puzzle_date),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
