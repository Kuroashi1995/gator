-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;


-- name: GetFeeds :many
SELECT f.name, f.url, u.name as user_name FROM feeds f
INNER JOIN users u
ON f.user_id = u.id;

-- name: GetFeedByUrl :one
SELECT *
FROM feeds f
WHERE f.url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds f
SET last_fetched_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE f.id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds f
JOIN users u
ON f.user_id = u.id
WHERE u.name = $1
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
