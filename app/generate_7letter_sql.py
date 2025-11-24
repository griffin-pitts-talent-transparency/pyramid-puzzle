#!/usr/bin/env python3
"""
generate_7letter_sql.py (12dicts / 6of12)
-----------------------------------------
1. Reads 6of12.txt (must exist in ./data)
2. Filters *7-letter lowercase words only*
3. Computes letter_mask and letters_unique
4. Generates a large SQL INSERT statement into:
       ./data/insert_7letter_words.sql
   (does NOT run the SQL against the database)
"""

import os
import re

# === Configuration ===
WORDLIST_TXT_PATH = os.path.join("data", "6of12.txt")

OUTPUT_SQL_PATH = os.path.join("data", "insert_7letter_words.sql")

TARGET_LEN = 7

def extract_words():
    if not os.path.exists(WORDLIST_TXT_PATH):
        raise SystemExit(f"❌ Word list not found at {WORDLIST_TXT_PATH}")

    print(f"📖 Reading words from {WORDLIST_TXT_PATH}...")

    with open(WORDLIST_TXT_PATH, "r", encoding="utf-8") as f:
        raw_words = [line.strip().lower() for line in f]

    # Only 7-letters, strictly a–z
    words = [
        w for w in raw_words
        if len(w) == TARGET_LEN and re.fullmatch(r"[a-z]+", w)
    ]

    print(f"✅ Found {len(words)} filtered English words (7 letters only)")
    return words


def word_to_mask(word: str) -> int:
    """Compute 26-bit mask of letters present in the word."""
    mask = 0
    for ch in set(word):
        mask |= 1 << (ord(ch) - ord("a"))
    return mask


def main():
    words = extract_words()

    print(f"🧮 Generating SQL file at {OUTPUT_SQL_PATH}...")

    values_rows = []
    for w in words:
        mask = word_to_mask(w)
        uniq = "".join(sorted(set(w)))

        # Format: ('word', length, letter_mask, 'letters_unique')
        row = f"('{w}', {TARGET_LEN}, {mask}, '{uniq}')"
        values_rows.append(row)

    if not values_rows:
        raise SystemExit("❌ No valid 7-letter words found. SQL file not generated.")

    # Build full SQL
    sql = (
        "INSERT INTO words (word, length, letter_mask, letters_unique)\n"
        "VALUES\n    "
        + ",\n    ".join(values_rows)
        + ";\n"
    )

    # Write SQL file
    with open(OUTPUT_SQL_PATH, "w", encoding="utf-8") as out:
        out.write(sql)

    print(f"✅ SQL generation complete. Total words: {len(values_rows)}")
    print(f"📄 SQL written to: {OUTPUT_SQL_PATH}")


if __name__ == "__main__":
    main()
