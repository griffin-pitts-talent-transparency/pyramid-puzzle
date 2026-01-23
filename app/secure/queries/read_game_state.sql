-- read_game_state.sql
SELECT gs.state
FROM game_state gs
JOIN users u ON u.id = gs.user_id
WHERE (
    (gs.date = :date)
    AND (:fingerprint IS NOT NULL AND trim(:fingerprint) != '' AND u.fingerprint = :fingerprint)
);
