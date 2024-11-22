-- name: CreateFeedFollow :many
WITH inserted_feed_follows AS (
    INSERT INTO feed_follows ( id, created_at, updated_at, user_id, feed_id)
    VALUES  (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)

SELECT
    inserted_feed_follows.*,
    feeds.name as feed_name,
    users.name as user_name
FROM inserted_feed_follows
INNER JOIN feeds
ON inserted_feed_follows.feed_id = feeds.id
INNER JOIN users
ON inserted_feed_follows.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT 
    feed_follows.*,
    feeds.name as feed_name
FROM feed_follows
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
USING feeds, users 
WHERE feed_follows.user_id = users.id
AND feed_follows.feed_id = feeds.id
AND users.id = $1
AND feeds.url = $2;

