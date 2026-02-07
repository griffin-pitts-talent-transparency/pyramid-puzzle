-- read_word_chain_game_state_by_sub.sql
SELECT guesses
FROM word_chain_game_state
JOIN users ON word_chain_game_state.user_id = users.id
WHERE users.cognito_sub = :cognito_sub
    AND puzzle_date = :puzzle_date;
