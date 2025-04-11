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
  is_mr_user = COALESCE($9, is_mr_user),
  is_mobile_user = COALESCE($10, is_mobile_user),
  goals = COALESCE($11, goals)
WHERE id = $1
RETURNING *;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

