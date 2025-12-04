-- read_solved_tier_stats.sql
WITH RECURSIVE
seq(n) AS (
    VALUES(0)
    UNION ALL
    SELECT n+1 FROM seq WHERE n < 63
),

players AS (
    SELECT COUNT(*) AS total_players
    FROM game_state
    WHERE (:date IS NULL OR game_state.date = :date)
      AND state IS NOT NULL
      AND trim(state) != ''
      AND json_array_length(json_extract(state, '$.guesses')) > 0
),

all_guesses AS (
    SELECT json_each.value AS guess
    FROM game_state, json_each(json_extract(state, '$.guesses'))
    WHERE (:date IS NULL OR game_state.date = :date)
),

expanded AS (
    SELECT
        guess,
        json_array_length(guess) AS len,
        CASE
            WHEN (
                SELECT MIN(
                    CASE
                        WHEN json_extract(guess, '$[' || seq.n || '].status') = 'match'
                        THEN 1 ELSE 0
                    END
                )
                FROM seq
                WHERE seq.n < json_array_length(guess)
            ) = 1
            THEN 1 ELSE 0
        END AS solved
    FROM all_guesses
)

SELECT
    -- SUM(CASE WHEN len = 4 AND solved = 1 THEN 1 ELSE 0 END) AS solved_4,
    -- SUM(CASE WHEN len = 5 AND solved = 1 THEN 1 ELSE 0 END) AS solved_5,
    -- SUM(CASE WHEN len = 6 AND solved = 1 THEN 1 ELSE 0 END) AS solved_6,
    -- SUM(CASE WHEN len = 7 AND solved = 1 THEN 1 ELSE 0 END) AS solved_7,

    COALESCE(ROUND(
        SUM(CASE WHEN len = 4 THEN solved ELSE 0 END) * 100.0 /
        (SELECT total_players FROM players), 2), 0
    ) AS pct_4,
    COALESCE(ROUND(
        SUM(CASE WHEN len = 5 THEN solved ELSE 0 END) * 100.0 /
        (SELECT total_players FROM players), 2), 0
    ) AS pct_5,
    COALESCE(ROUND(
        SUM(CASE WHEN len = 6 THEN solved ELSE 0 END) * 100.0 /
        (SELECT total_players FROM players), 2), 0
    ) AS pct_6,
    COALESCE(ROUND(
        SUM(CASE WHEN len = 7 THEN solved ELSE 0 END) * 100.0 /
        (SELECT total_players FROM players), 2), 0
    ) AS pct_7
FROM expanded;
