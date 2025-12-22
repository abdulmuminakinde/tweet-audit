package cmd

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"time"

	"github.com/abdulmuminakinde/tweet-audit/internal/config"
	"github.com/abdulmuminakinde/tweet-audit/internal/core"
	"github.com/abdulmuminakinde/tweet-audit/internal/database"
)

const batchSize = 500

type State struct {
	Config  *config.Config
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

func (s *State) StreamTweets(ctx context.Context, file io.Reader) error {
	r := bufio.NewReader(file)
	dec := json.NewDecoder(r)

	batch := make([]InsertTweetParams, 0, batchSize)

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

		full_text := normalizeText(tweet.Tweet.FullText)

		url, err := getTweetUrl(tweet, s.Config.Username)
		if err != nil {
			return err
		}

		batch = append(batch, InsertTweetParams{
			TweetID:           tweet.Tweet.ID,
			CreatedAt:         parsedCreatedAt,
			FullText:          full_text,
			Retweeted:         tweet.Tweet.Retweeted,
			PossiblySensitive: tweet.Tweet.PossiblySensitive,
			Url:               url,
		})

		if len(batch) >= batchSize {
			if err := s.batchInsertTweets(ctx, batch); err != nil {
				return err
			}
			batch = batch[:0]
		}

	}

	if len(batch) > 0 {
		if err := s.batchInsertTweets(ctx, batch); err != nil {
			return err
		}
	}

	return nil
}

func (s *State) batchInsertTweets(ctx context.Context, tweets []InsertTweetParams) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Queries.WithTx(tx)

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
