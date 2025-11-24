-- read_game_state.sql
SELECT state
FROM game_state
WHERE (
    (user_id = (SELECT id FROM users WHERE fingerprint = :fingerprint))
    AND (date = :date)
);
    
