DROP TABLE IF EXISTS words;

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS words (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    word           TEXT NOT NULL UNIQUE,
    length         INTEGER NOT NULL,
    letter_mask    INTEGER NOT NULL,
    letters_unique TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_words_length ON words(length);
CREATE INDEX IF NOT EXISTS idx_words_mask   ON words(letter_mask);
