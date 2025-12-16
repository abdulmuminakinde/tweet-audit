package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/abdulmuminakinde/tweet-audit/internal/core"
)

func StreamTweets(file io.Reader) error {
	r := bufio.NewReader(file)
	dec := json.NewDecoder(r)

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
		var tweet core.Post
		if err := dec.Decode(&tweet); err != nil {
			return err
		}

		url, err := getTweetUrl(tweet)
		if err != nil {
			return err
		}

		fmt.Println(tweet.Tweet.FullText, url)
		fmt.Println("===================")
	}

	return nil

}

func getTweetUrl(tweetObject core.Post) (string, error) {
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
