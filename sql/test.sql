SELECT *
FROM users
WHERE id IN (
    SELECT
        CASE
            WHEN receiver_id = 10 THEN sender_id
            ELSE receiver_id
        END
    FROM messages
    WHERE sender_id = 10
       OR receiver_id = 10
    )