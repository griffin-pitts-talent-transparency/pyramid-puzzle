-- initialize_game_state.sql
INSERT INTO game_state (user_id, state, date)
VALUES (
    (SELECT id FROM users WHERE fingerprint = :fingerprint),
    json_object('guesses', json_array()),
    :date
);
