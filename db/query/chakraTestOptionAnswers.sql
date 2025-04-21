-- name: CreateChakraTestOptionAnswersBatch :many
INSERT INTO chakra_test_option_answers (
    unique_id,
    email,
    unique_code,
    question,
    answer
)
SELECT unnest($1::uuid[]), unnest($2::text[]), unnest($3::text[]), unnest($4::text[]), unnest($5::text[])
RETURNING *; 