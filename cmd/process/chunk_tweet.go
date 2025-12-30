package main

import "github.com/abdulmuminakinde/tweet-audit/internal/database"

func chunkTweets(tweets []database.GetTweetsRow, size int) [][]database.GetTweetsRow {
	var chunks [][]database.GetTweetsRow

	for i := 0; i < len(tweets); i += size {
		end := i + size
		if end > len(tweets) {
			end = len(tweets)
		}
		chunks = append(chunks, tweets[i:end])
	}

	return chunks
}
