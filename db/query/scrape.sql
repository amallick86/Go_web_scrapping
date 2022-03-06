-- name: CreateScrape :one
INSERT INTO scrape (
  user_id, url,scrapped
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetOwnScrape :many
SELECT scrape.id,scrape.url,scrape.scrapped,scrape.created_at, users.username FROM scrape INNER JOIN users ON scrape.user_id = users.id
WHERE scrape.user_id = $1 AND scrape.id <= $2 LIMIT 10;

-- name: CountOwnScrape :one
SELECT COUNT(*) FROM scrape
WHERE user_id = $1;

-- name: GetScrape :many
SELECT scrape.id,scrape.url,scrape.scrapped,scrape.created_at, users.username  FROM scrape INNER JOIN users ON scrape.user_id = users.id
WHERE  scrape.id <= $1 LIMIT 10;

-- name: CountScrape :one
SELECT COUNT(*) FROM scrape;

-- name: Search :many
SELECT scrape.id,scrape.url,scrape.scrapped,scrape.created_at, users.username FROM scrape INNER JOIN users ON scrape.user_id = users.id WHERE url @@ $1 ;

-- name: Filter :many
SELECT scrape.id,scrape.url,scrape.scrapped,scrape.created_at, users.username FROM scrape INNER JOIN users ON scrape.user_id = users.id WHERE scrape.created_at  BETWEEN $1 AND $2;

-- name: MinDate :one
SELECT created_at FROM scrape WHERE created_at=( SELECT MIN(created_at) FROM scrape ) ;