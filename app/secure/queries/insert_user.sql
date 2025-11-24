-- insert_user.sql
INSERT INTO users (
    fingerprint,
    ip,
    user_agent,
    device_type,
    os
)
SELECT
    :fingerprint,
    :ip,
    :user_agent,
    :device_type,
    :os
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE fingerprint = :fingerprint
)
RETURNING id;
