-- initialize_word_chain_game_state_by_sub.sql
INSERT INTO word_chain_game_state (user_id, puzzle_date, guesses)
SELECT
    users.id,
    :puzzle_date,
    json('[]')
FROM users
WHERE users.cognito_sub = :cognito_sub
AND NOT EXISTS (
    SELECT 1 FROM word_chain_game_state
    WHERE word_chain_game_state.user_id = users.id
    AND word_chain_game_state.puzzle_date = :puzzle_date
);
