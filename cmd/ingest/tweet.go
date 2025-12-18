package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/abdulmuminakinde/tweet-audit/internal/core"
)

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
