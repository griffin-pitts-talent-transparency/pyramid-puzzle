-- insert_login_stream.sql
INSERT INTO login_streaks (user_id, login_date, streak)
SELECT
    :user_id,
    DATE('now', 'localtime'),
    COALESCE((
        SELECT streak + 1
        FROM login_streaks
        WHERE user_id = :user_id AND login_date = DATE('now', 'localtime', '-1 day')
    ), 1)
WHERE NOT EXISTS (
    SELECT 1
    FROM login_streaks
    WHERE user_id = :user_id AND login_date = DATE('now', 'localtime')
);
