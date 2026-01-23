-- read_user_by_sub.sql
SELECT
    id,
    cognito_sub,
    cognito_username,
    fingerprint,
    ip,
    user_agent,
    device_type,
    os,
    creation_date
FROM users
WHERE (
    (:cognito_sub IS NOT NULL AND TRIM(:cognito_sub) != '')
    AND (cognito_sub = :cognito_sub)
)
LIMIT 1;
