-- deploy_word_chain_schema.sql
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS word_chain_puzzles (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    puzzle_date  DATE UNIQUE NOT NULL,
    word_chain   TEXT NOT NULL
);
