-- name: CreateUserProfileMr :exec
INSERT INTO users_profile_mr (
  user_id
) VALUES (
  $1
)
RETURNING *;

-- name: GetUserProfileMrByUserId :one
SELECT * FROM users_profile_mr
WHERE user_id = $1;

-- name: GetUserProfileMrTime :one
SELECT total_time_spent_in_mins FROM users_profile_mr
WHERE user_id = $1;