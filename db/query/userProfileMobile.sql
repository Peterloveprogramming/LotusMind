-- name: CreateUserProfileMobile :exec
INSERT INTO users_profile_mobile (
  user_id
) VALUES (
  $1
)
RETURNING *;

-- name: GetUserProfileMobileByUserId :one
SELECT * FROM users_profile_mobile
WHERE user_id = $1;

-- name: GetUserProfileMobileTime :one
SELECT total_time_spent_in_mins FROM users_profile_mobile
WHERE user_id = $1;