-- read_solved_games_all_time.sql
WITH RECURSIVE
seq(n) AS (
    VALUES(0)
    UNION ALL
    SELECT n+1 FROM seq WHERE n < 6
),

valid_users AS (
    SELECT id
    FROM users
    WHERE cognito_sub IS NOT NULL AND trim(cognito_sub) != ''
),

all_guesses AS (
    SELECT
        g.user_id,
        g.date,
        json_each.value AS guess
    FROM game_state g
    JOIN valid_users u ON g.user_id = u.id,
         json_each(json_extract(g.state, '$.guesses'))
),

guess_status AS (
    SELECT
        user_id,
        date,
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
    FROM all_guesses, seq
    GROUP BY user_id, date, guess
),

solved_per_day AS (
    SELECT
        user_id,
        date,
        len,
        MAX(solved) AS solved
    FROM guess_status
    GROUP BY user_id, date, len
),

daily_fully_solved AS (
    SELECT
        user_id,
        date
    FROM solved_per_day
    WHERE solved = 1
    GROUP BY user_id, date
    HAVING COUNT(DISTINCT len) = 4
)

SELECT
    u.cognito_username,
    COUNT(*) AS total_solved_games
FROM daily_fully_solved d
JOIN users u ON u.id = d.user_id
GROUP BY d.user_id
ORDER BY total_solved_games DESC
LIMIT 100;
