SELECT COUNT(*) AS incorrect_guess_count
FROM (
    SELECT guess.value AS guess_array
    FROM game_state,
         json_each(game_state.state, '$.guesses') AS guess
    WHERE user_id = (SELECT id FROM users WHERE fingerprint = :fingerprint)
      AND date = :date
)
WHERE EXISTS (
    SELECT 1
    FROM json_each(guess_array)
    WHERE json_extract(json_each.value, '$.status') != 'match'
);
