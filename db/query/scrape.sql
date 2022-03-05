-- name: CreateScrape :one
INSERT INTO scrape (
  user_id, url,scrapped
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetOwnScrape :many
SELECT * FROM scrape
WHERE user_id = $1 AND id <= $2 LIMIT 5;

-- name: CountOwnScrape :one
SELECT COUNT(*) FROM scrape
WHERE user_id = $1;

-- name: GetScrape :many
SELECT * FROM scrape
WHERE  id <= $1 LIMIT 5;

-- name: CountScrape :one
SELECT COUNT(*) FROM scrape;

-- name: Search :many
SELECT * FROM scrape WHERE url @@ $1 ;

-- name: Filter :many
SELECT * FROM scrape WHERE created_at  BETWEEN $1 AND $2;