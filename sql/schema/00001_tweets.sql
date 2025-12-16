-- +goose Up
-- +goose StatementBegin
CREATE TABLE tweets
(
    tweet_id TEXT PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL, -- time tweet was created
    full_text TEXT NOT NULL,
    possibly_sensitive BOOLEAN NOT NULL,
    retweeted BOOLEAN NOT NULL
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE tweets;
-- +goose StatementEnd


