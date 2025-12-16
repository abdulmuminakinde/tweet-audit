-- name: InsertTweet :exec
INSERT INTO tweets
(
  tweet_id,
  created_at,
  full_text,
  possibly_sensitive,
  retweeted
) VALUES (
  $1, $2, $3, $4, $5
  )
RETURNING *;

