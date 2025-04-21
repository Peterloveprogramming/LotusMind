-- name: GetChakraTestResults :many
SELECT unique_id, email, chakra_name, chakra_score, chakra_status, created_at
FROM chakra_test_results
WHERE email = $1 AND deleted_at = '0001-01-01 00:00:00Z'
ORDER BY created_at DESC; 


 -- name: CreateChakraTestResult :one
INSERT INTO chakra_test_results (
    unique_id,
    email,
    chakra_name,
    chakra_score,
    chakra_status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *; 
