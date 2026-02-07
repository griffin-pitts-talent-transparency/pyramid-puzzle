-- read_login_streak.sql
SELECT
    ls.streak
FROM login_streaks AS ls 
WHERE (
    (:user_id IS NOT NULL AND TRIM(:user_id) != '' AND ls.user_id = :user_id)
    AND (ls.login_date = DATE('now', 'localtime'))
)
LIMIT 1;