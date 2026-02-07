-- read_user.sql
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
    (:fingerprint IS NOT NULL AND TRIM(:fingerprint) != '')
    AND (fingerprint = :fingerprint)
)
LIMIT 1;
