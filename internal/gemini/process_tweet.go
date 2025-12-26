package gemini

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abdulmuminakinde/tweet-audit/internal/database"
)

func ProcessBatchedTweets(ctx context.Context, client *Client, tweets []database.GetTweetsRow) error {
	const batchSize = 40

	batches := chunkTweets(tweets, batchSize)

	var allResults []TweetAnalysisResult

	for i, batch := range batches {
		fmt.Printf("Processing batch %d/%d (%d tweets)...\n", i+1, len(batches), len(batch))

		results, err := client.AnalyzeTweets(ctx, batch)
		if err != nil {
			log.Printf("Batch %d failed: %v", i+1, err)
		}

		allResults = append(allResults, results...)

		if i < len(batches)-1 {
			time.Sleep(2 * time.Second)
		}
	}
	fmt.Println(allResults[:30])
	return nil
}
