package main

import (
	"log"

	cmd "github.com/abdulmuminakinde/tweet-audit/cmd/ingest"
	"github.com/abdulmuminakinde/tweet-audit/internal/core"
)

func main() {
	file, err := core.LoadTweets("./internal/tweets.js")
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.StreamTweets(file)
	if err != nil {
		log.Fatal(err)
	}

}
