-- name: CreateTibetanSingingBowlMr :one
INSERT INTO tibetan_singing_bowl_mr (
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

-- name: UpdateTibetanSingingBowlMrByUniqueID :one
UPDATE tibetan_singing_bowl_mr
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


-- name: GetTibetanSingingBowlMrByUniqueID :one
SELECT * 
FROM tibetan_singing_bowl_mr 
WHERE unique_id = $1;

-- name: GetTibetanSingingBowlMrByUuid :one
SELECT * 
FROM tibetan_singing_bowl_mr 
WHERE uuid = $1;

-- name: GetTibetanSingingBowlMrByUserID :many
SELECT 
  tbm.*
FROM 
  tibetan_singing_bowl_mr AS tbm
JOIN 
  session_logs AS sl ON tbm.uuid = sl.uuid
WHERE 
  sl.user_id = $1;



-- name: DeleteTibetanSingingBowlMrByUniqueID :one
DELETE FROM tibetan_singing_bowl_mr 
WHERE unique_id = $1
RETURNING *;
