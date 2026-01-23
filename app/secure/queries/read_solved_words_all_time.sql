-- read_solved_words_all_time.sql
WITH RECURSIVE
seq(n) AS (
    VALUES(0)
    UNION ALL
    SELECT n+1 FROM seq WHERE n < 63
),

valid_users AS (
    SELECT id
    FROM users
    WHERE cognito_sub IS NOT NULL AND trim(cognito_sub) != ''
),

user_guesses AS (
    SELECT
        game_state.user_id,
        json_each.value AS guess
    FROM game_state
    JOIN valid_users ON game_state.user_id = valid_users.id,
         json_each(json_extract(game_state.state, '$.guesses'))
),

expanded AS (
    SELECT
        user_id,
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
    FROM user_guesses, seq
    GROUP BY user_id, guess
),

aggregated AS (
    SELECT
        user_id,
        SUM(solved) AS total_solved
    FROM expanded
    GROUP BY user_id
)

SELECT
    users.cognito_username,
    aggregated.total_solved
FROM aggregated
JOIN users ON users.id = aggregated.user_id
ORDER BY aggregated.total_solved DESC
LIMIT 100;
