-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (
    id,
    created_at,
    updated_at,
    user_id,
    feed_id
    ) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
    )
    RETURNING *
)
SELECT iff.*, f.name as feed_name, u.name as user_name
FROM inserted_feed_follow iff
JOIN feeds f
ON f.id = iff.feed_id
JOIN users u
ON u.id = iff.user_id;


-- name: GetFeedFollowsForUser :many
SELECT ff.*, f.name as feed_name, u.name as user_name
FROM feed_follows ff
JOIN feeds f
ON f.id = ff.feed_id
JOIN users u
ON u.id = ff.user_id
WHERE u.name = $1;

-- name: DeleteFeedFollowByUserFeed :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;
