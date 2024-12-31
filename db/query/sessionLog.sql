-- name: CreateSessionLog :one
INSERT INTO session_logs (
  user_id,
  session_type,
  session_platform
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetSessionLogByUUID :one
SELECT * FROM session_logs
WHERE uuid = $1;

-- name: GetUserIdFromSessionLogUuid :one
SELECT user_id FROM session_logs
WHERE uuid = $1;

-- name: GetSessionsLogByUserIdWithOffset :many
SELECT * FROM session_logs
WHERE user_id = $1
LIMIT $2
OFFSET $3;

-- name: GetSessionsLogByUserId :many
SELECT * FROM session_logs
WHERE user_id = $1;

-- name: DeleteSessionLog :exec
DELETE FROM session_logs
WHERE uuid = $1;

