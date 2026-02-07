-- insert_user.sql
INSERT INTO users (
    fingerprint,
    cognito_sub,
    cognito_username,
    ip,
    user_agent,
    device_type,
    os
)
SELECT
    :fingerprint,
    :cognito_sub,
    :cognito_username,
    :ip,
    :user_agent,
    :device_type,
    :os
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE fingerprint = :fingerprint
)
RETURNING id;
