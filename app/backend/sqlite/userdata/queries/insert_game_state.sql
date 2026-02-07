-- insert_game_state.sql
INSERT INTO game_state (user_id, state)
VALUES (
    (SELECT id FROM users WHERE fingerprint = :fingerprint),
    :state
);
