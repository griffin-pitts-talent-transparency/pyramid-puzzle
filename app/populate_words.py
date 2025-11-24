#!/usr/bin/env python3
"""
populate_words.py (12dicts / 2of12inf, local file only)
-------------------------------------------------------
1. Assumes 2of12inf.txt from the 12dicts project is already present in ./data
2. Reads 2of12inf.txt
3. Filters 3–6 letter lowercase words (no abbreviations, no punctuation)
4. Computes letter_mask and letters_unique
5. Inserts into pyramid_puzzle.db
"""

import os
import re
import sqlite3

# === Configuration ===
DB_PATH = os.path.join("data", "pyramid_puzzle.db")

# Manually downloaded:
WORDLIST_TXT_PATH = os.path.join("data", "6of12.txt")

MIN_LEN, MAX_LEN = 3, 6


def extract_words():
    if not os.path.exists(WORDLIST_TXT_PATH):
        raise SystemExit(
            f"❌ Word list not found at {WORDLIST_TXT_PATH}"
        )

    print(f"📖 Reading words from {WORDLIST_TXT_PATH}...")

    with open(WORDLIST_TXT_PATH, "r", encoding="utf-8") as f:
        raw_words = [line.strip().lower() for line in f]

    # 3–6 letters, strictly a–z. 2of12inf already excludes abbreviations,
    # caps, hyphenations, etc, but this keeps things extra clean.
    words = [
        w for w in raw_words
        if MIN_LEN <= len(w) <= MAX_LEN and re.fullmatch(r"[a-z]+", w)
    ]

    print(f"✅ Found {len(words)} filtered English words (3-6 letters)")
    return words


def word_to_mask(word: str) -> int:
    """Compute 26-bit mask of letters present in the word."""
    mask = 0
    for ch in set(word):
        mask |= 1 << (ord(ch) - ord("a"))
    return mask


def main():
    if not os.path.exists(DB_PATH):
        raise SystemExit(f"❌ Database not found at {DB_PATH}")

    words = extract_words()

    print(f"🧮 Inserting into {DB_PATH}...")

    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()

    insert_sql = """
        INSERT OR IGNORE INTO words (word, length, letter_mask, letters_unique)
        VALUES (?, ?, ?, ?);
    """

    batch = []
    for w in words:
        mask = word_to_mask(w)
        uniq = "".join(sorted(set(w)))
        batch.append((w, len(w), mask, uniq))

    cur.executemany(insert_sql, batch)
    conn.commit()
    conn.close()

    print(f"✅ Inserted {len(batch)} clean English words.")


if __name__ == "__main__":
    main()
