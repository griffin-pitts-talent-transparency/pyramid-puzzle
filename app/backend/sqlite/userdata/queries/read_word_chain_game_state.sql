-- read_word_chain_game_state
SELECT
    json_extract(guess.value, '$.word') AS word,
    json_extract(guess.value, '$.result') AS result
FROM word_chain_game_state,
    json_each(guesses) AS guess
WHERE user_id = :user_id AND puzzle_date = :puzzle_date;
