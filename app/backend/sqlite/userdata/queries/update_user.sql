-- update_user.sql
UPDATE users
SET
    cognito_sub = CASE
        WHEN :cognito_sub IS NOT NULL AND TRIM(:cognito_sub) != '' AND (cognito_sub IS NULL OR TRIM(cognito_sub) == '')
        THEN :cognito_sub
        ELSE cognito_sub
    END,
    cognito_username = CASE
        WHEN :cognito_username IS NOT NULL AND TRIM(:cognito_username) != ''
             AND (cognito_username IS NULL OR TRIM(cognito_username) == '' OR :cognito_username != cognito_username)
        THEN :cognito_username
        ELSE cognito_username
    END
WHERE id = :user_id;
