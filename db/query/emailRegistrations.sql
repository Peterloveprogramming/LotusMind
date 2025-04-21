-- name: CreateUserEmail :one
INSERT INTO email_registrations  (
  unique_id,
  email,
  chakra_info,
  language,
  unique_code,
  ip,
  country
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING unique_id, email, chakra_info, language, unique_code, ip, country, created_at, deleted_at;


-- name: GetLatestEmailRegistration :one
SELECT unique_id, email, unique_code,created_at, deleted_at FROM email_registrations
WHERE email = $1
ORDER BY created_at DESC
LIMIT 1;


-- name: GetEmailRegistrationByTestNum :one
SELECT unique_id, email, language, chakra_info, unique_code, created_at, deleted_at
FROM email_registrations
WHERE email = $1
ORDER BY created_at
OFFSET $2 LIMIT 1;

-- name: GetReportByCode :one
SELECT unique_id, email, chakra_info, language, unique_code, created_at, deleted_at
FROM email_registrations
WHERE unique_code = $1 AND deleted_at = '0001-01-01 00:00:00+00'
LIMIT 1;