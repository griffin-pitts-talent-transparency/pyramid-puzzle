UPDATE game_state
SET state = json_set(
    state,
    '$.guesses[' || COALESCE(json_array_length(json_extract(state, '$.guesses')), 0) || ']',
    json(:guess)
)
WHERE user_id = (SELECT id FROM users WHERE fingerprint = :fingerprint)
  AND date = :date;
