-- name: EnqueueCommand :exec
INSERT INTO message_queue (command_type, payload)
VALUES ($1, $2);

-- name: DequeueCommand :one
DELETE FROM message_queue 
WHERE id = (
    SELECT id FROM message_queue 
    ORDER BY created_at ASC 
    FOR UPDATE SKIP LOCKED 
    LIMIT 1
)
RETURNING *;
