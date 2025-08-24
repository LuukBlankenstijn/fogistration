-- name: GetClientById :one
SELECT * FROM clients
WHERE id = $1;

-- name: GetAllClients :many
SELECT * FROM clients;

-- name: UpdateClientLastSeen :exec
UPDATE clients
SET last_seen = NOW()
WHERE ip = $1;

-- name: UpsertClient :one
INSERT INTO clients (
    ip
) VALUES (
    $1
)
ON CONFLICT (ip)
DO UPDATE SET
    last_seen = NOW()
RETURNING *;
