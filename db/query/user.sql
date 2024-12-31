-- name: CreateUser :one
INSERT INTO users (
  email,
  first_name,
  last_name,
  gender,
  birth_date,
  country,
  hashed_password,
  is_mr_user,
  is_mobile_user,
  goals
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9,$10
)
RETURNING *;

-- name: GetUsersByCountry :many
SELECT * FROM users
where country = $1
ORDER BY created_at;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

