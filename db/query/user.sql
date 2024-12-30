-- name: CreateUser :one
INSERT INTO users (
  email,
  first_name,
  last_name,
  gender,
  birth_date,
  country,
  hashed_password,
  goals,
  platform
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetUsersByPlatform :many
SELECT * FROM users
where platform = $1
ORDER BY created_at;

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

-- name: UpdateUser :one
UPDATE users
SET 
    email = COALESCE($2, email), 
    first_name = COALESCE($3, first_name), 
    last_name = COALESCE($4, last_name),
    gender = COALESCE($5, gender), 
    birth_date = COALESCE($6, birth_date), 
    country = COALESCE($7, country), 
    hashed_password = COALESCE($8, hashed_password), 
    goals = COALESCE($9, goals), 
    platform = COALESCE($10, platform), 
    password_changed_at = COALESCE($11, password_changed_at)
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

