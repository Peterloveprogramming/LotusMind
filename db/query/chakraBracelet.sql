-- name: GetChakraBracelet :many
WITH custom_bracelets AS (
    SELECT chakra, name, image_url, product_link, type, created_at, deleted_at
    FROM chakra_bracelet 
    WHERE chakra = ANY(sqlc.arg(chakras)::text[]) AND type = 0 -- Use sqlc.arg and type cast
    ORDER BY created_at DESC
),
random_bracelets AS (
    SELECT chakra, name, image_url, product_link, type, created_at, deleted_at
    FROM chakra_bracelet 
    WHERE type = 1
    ORDER BY RANDOM()
    LIMIT 2
)
SELECT * FROM random_bracelets	
UNION ALL
SELECT * FROM custom_bracelets;
