package cmd

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/abdulmuminakinde/tweet-audit/internal/core"
	"github.com/abdulmuminakinde/tweet-audit/internal/database"
)

const BatchSize = 500

type Config struct {
	DB      *sql.DB
	Queries *database.Queries
}

type InsertTweetParams struct {
	TweetID           string
	CreatedAt         time.Time
	FullText          string
	PossiblySensitive bool
	Retweeted         bool
	Url               string
}

func (c *Config) StreamTweets(ctx context.Context, file io.Reader) error {
	r := bufio.NewReader(file)
	dec := json.NewDecoder(r)

	batch := make([]InsertTweetParams, 0, BatchSize)

	// skip text before the actual JSON array
	// aeems to me like I could just edit the
	// js file into json though
	for {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}
		if b == '[' {
			// push back "[" into the reader
			r.UnreadByte()
			break
		}
	}

	if _, err := dec.Token(); err != nil {
		return err
	}

	for dec.More() {
		var tweet core.RawTweet
		if err := dec.Decode(&tweet); err != nil {
			return err
		}

		parsedCreatedAt, err := parseTweetDate(tweet.Tweet.CreatedAt)
		if err != nil {
			return err
		}

		url, err := getTweetUrl(tweet)
		if err != nil {
			return err
		}

		batch = append(batch, InsertTweetParams{
			TweetID:           tweet.Tweet.ID,
			CreatedAt:         parsedCreatedAt,
			FullText:          tweet.Tweet.FullText,
			Retweeted:         tweet.Tweet.Retweeted,
			PossiblySensitive: tweet.Tweet.PossiblySensitive,
			Url:               url,
		})

		if len(batch) >= BatchSize {
			if err := c.batchInsertTweets(ctx, batch); err != nil {
				return err
			}
			batch = batch[:0]
		}

	}

	if len(batch) > 0 {
		if err := c.batchInsertTweets(ctx, batch); err != nil {
			return err
		}
	}

	return nil
}

func getTweetUrl(tweetObject core.RawTweet) (string, error) {
	var url string
	var username string = "lanrey_waju"
	var tweetID string

	tweetID = tweetObject.Tweet.ID
	if tweetID == "" {
		return "", errors.New("unable to get url for tweet")
	}

	url = fmt.Sprintf("https://x.com/%s/status/%s", username, tweetID)
	return url, nil

}

func parseTweetDate(createdAt string) (time.Time, error) {
	// Twitter's time format is exactly time.RubyDate (Mon Sep 20 20:16:27 +0000 2025)
	timeFormat := time.RubyDate

	parsedTime, err := time.Parse(timeFormat, createdAt)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date for tweet created_at %q: %w", createdAt, err)
	}

	return parsedTime, nil
}

func (c *Config) batchInsertTweets(ctx context.Context, tweets []InsertTweetParams) error {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := c.Queries.WithTx(tx)

	for _, tweet := range tweets {
		if err := qtx.InsertTweet(ctx, database.InsertTweetParams{
			TweetID:           tweet.TweetID,
			CreatedAt:         tweet.CreatedAt,
			FullText:          tweet.FullText,
			PossiblySensitive: tweet.PossiblySensitive,
			Retweeted:         tweet.Retweeted,
			Url:               tweet.Url,
		}); err != nil {
			return err
		}
	}

	return tx.Commit()
}
