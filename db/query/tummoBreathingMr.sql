-- name: CreateTummoBreathingMr :one
INSERT INTO tummo_breathing_mr (
  uuid,
  start_mood_rating,
  start_mood,
  finish_mood_rating,
  finish_mood,
  session_completed
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateTummoBreathingMrByUniqueID :one
UPDATE tummo_breathing_mr
SET
  uuid = COALESCE($2, uuid),
  start_mood_rating = COALESCE($3, start_mood_rating),
  start_mood = COALESCE($4, start_mood),
  finish_mood_rating = COALESCE($5, finish_mood_rating),
  finish_mood = COALESCE($6, finish_mood),
  session_completed = COALESCE($7, session_completed),
  started_at = COALESCE($8, started_at),
  ends_at = COALESCE($9, ends_at),
  deleted_at = COALESCE($10, deleted_at)
WHERE unique_id = $1
RETURNING *;

-- name: UpdateTummoBreathingMrByUuid :one
UPDATE tummo_breathing_mr
SET
  start_mood_rating = COALESCE($2, start_mood_rating),
  start_mood = COALESCE($3, start_mood),
  finish_mood_rating = COALESCE($4, finish_mood_rating),
  finish_mood = COALESCE($5, finish_mood),
  session_completed = COALESCE($6, session_completed),
  started_at = COALESCE($7, started_at),
  ends_at = COALESCE($8, ends_at),
  deleted_at = COALESCE($9, deleted_at)
WHERE uuid = $1
RETURNING *;

-- name: UpdateTummoBreathingMrStartingMoodByUuid :one
UPDATE tummo_breathing_mr
SET
  start_mood_rating = COALESCE($2, start_mood_rating),
  start_mood = COALESCE($3, start_mood)
WHERE uuid = $1
RETURNING *;

-- name: UpdateTummoBreathingMrFinishingMoodByUuid :one
UPDATE tummo_breathing_mr
SET
 finish_mood_rating = COALESCE($2, finish_mood_rating),
  finish_mood = COALESCE($3, finish_mood),
  session_completed = COALESCE($4, session_completed),
  ends_at = COALESCE($5, ends_at)
WHERE uuid = $1
RETURNING *;

-- name: UpdateTummoBreathingMrQuitByUuid :one
UPDATE tummo_breathing_mr
SET
  ends_at = COALESCE($2, ends_at)
WHERE uuid = $1
RETURNING *;

-- name: GetTummoBreathingMrByUniqueID :one
SELECT * 
FROM tummo_breathing_mr 
WHERE unique_id = $1;

-- name: GetTummoBreathingMrByUuid :one
SELECT * 
FROM tummo_breathing_mr 
WHERE uuid = $1;

-- name: GetTummoBreathingMrByUserID :many
SELECT 
  tbm.*
FROM 
  tummo_breathing_mr AS tbm
JOIN 
  session_logs AS sl ON tbm.uuid = sl.uuid
WHERE 
  sl.user_id = $1;

-- name: SoftDeleteTummoBreathingMrByUniqueID :one
UPDATE tummo_breathing_mr 
SET deleted_at = NOW()
WHERE unique_id = $1 AND (deleted_at IS NULL OR deleted_at = '0001-01-01 00:00:00Z')
RETURNING *;