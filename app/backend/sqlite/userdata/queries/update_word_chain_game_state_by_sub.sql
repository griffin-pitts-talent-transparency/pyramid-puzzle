-- update_word_chain_game_state_by_sub.sql
UPDATE word_chain_game_state
SET guesses = json_insert(
    guesses,
    '$[' || json_array_length(guesses) || ']',
    json_object(
        'word',   :word,
        'result', :result
    )
)
WHERE user_id = (
    SELECT id FROM users WHERE cognito_sub = :cognito_sub LIMIT 1
)
AND puzzle_date = :puzzle_date;
