-- select_word.sql
SELECT word, length, letter_mask
FROM words
WHERE length = :length
  AND (
    -- allow first word to ignore bitcount (prev_mask = 0, required_shared = 0)
    ( :required_shared = 0 AND :prev_mask = 0 )
    OR (
        -- otherwise calculate bitcount of shared bits:
        (
          ((letter_mask & :prev_mask) & 1)
          + (((letter_mask & :prev_mask) >> 1) & 1)
          + (((letter_mask & :prev_mask) >> 2) & 1)
          + (((letter_mask & :prev_mask) >> 3) & 1)
          + (((letter_mask & :prev_mask) >> 4) & 1)
          + (((letter_mask & :prev_mask) >> 5) & 1)
          + (((letter_mask & :prev_mask) >> 6) & 1)
          + (((letter_mask & :prev_mask) >> 7) & 1)
          + (((letter_mask & :prev_mask) >> 8) & 1)
          + (((letter_mask & :prev_mask) >> 9) & 1)
          + (((letter_mask & :prev_mask) >>10) & 1)
          + (((letter_mask & :prev_mask) >>11) & 1)
          + (((letter_mask & :prev_mask) >>12) & 1)
          + (((letter_mask & :prev_mask) >>13) & 1)
          + (((letter_mask & :prev_mask) >>14) & 1)
          + (((letter_mask & :prev_mask) >>15) & 1)
          + (((letter_mask & :prev_mask) >>16) & 1)
          + (((letter_mask & :prev_mask) >>17) & 1)
          + (((letter_mask & :prev_mask) >>18) & 1)
          + (((letter_mask & :prev_mask) >>19) & 1)
          + (((letter_mask & :prev_mask) >>20) & 1)
          + (((letter_mask & :prev_mask) >>21) & 1)
          + (((letter_mask & :prev_mask) >>22) & 1)
          + (((letter_mask & :prev_mask) >>23) & 1)
          + (((letter_mask & :prev_mask) >>24) & 1)
          + (((letter_mask & :prev_mask) >>25) & 1)
        ) = :required_shared
    )
)
ORDER BY id
LIMIT 1 OFFSET :offset;
