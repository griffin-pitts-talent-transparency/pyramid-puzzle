-- read_total_guess_count.sql
SELECT COUNT(*) AS total_guess_count
FROM game_state, json_each(json_extract(game_state.state, '$.guesses'));
