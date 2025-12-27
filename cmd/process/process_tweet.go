package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/abdulmuminakinde/tweet-audit/internal/database"
	"github.com/abdulmuminakinde/tweet-audit/internal/gemini"
)

type TweetProcessor struct {
	client    *gemini.Client
	dbQueries *database.Queries
}

func NewTweetProcessor(client *gemini.Client, queries *database.Queries) *TweetProcessor {
	return &TweetProcessor{
		client:    client,
		dbQueries: queries,
	}
}

func (p *TweetProcessor) ProcessAllTweets(ctx context.Context) error {
	tweets, err := p.dbQueries.GetTweets(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch tweets: %w", err)
	}

	const batchSize = 40

	batches := chunkTweets(tweets, batchSize)

	results := p.ProcessBatchesConcurrently(ctx, batches)

	log.Printf("Saving %d results...", len(results))
	return p.saveResults(ctx, results)
}

func (p *TweetProcessor) ProcessBatchesConcurrently(ctx context.Context, batches [][]database.GetTweetsRow) []gemini.TweetAnalysisResult {
	const numWorkers = 3

	jobs := make(chan []database.GetTweetsRow, len(batches))
	results := make(chan []gemini.TweetAnalysisResult, len(batches))

	limiter := time.NewTicker(6 * time.Second)
	defer limiter.Stop()

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go p.worker(ctx, i+1, jobs, results, limiter, &wg)
	}

	for _, batch := range batches {
		jobs <- batch
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var allResults []gemini.TweetAnalysisResult
	for batchResults := range results {
		allResults = append(allResults, batchResults...)
	}

	return allResults
}

func (p *TweetProcessor) worker(
	ctx context.Context,
	workerID int,
	jobs <-chan []database.GetTweetsRow,
	results chan<- []gemini.TweetAnalysisResult,
	limiter *time.Ticker,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for batch := range jobs {
		<-limiter.C

		log.Printf("Worker %d processing batch of %d tweets...", workerID, len(batch))

		batchResults, err := p.client.AnalyzeTweets(ctx, batch)
		if err != nil {
			log.Printf("Worker %d: batch failed: %v", workerID, err)
			continue
		}

		log.Printf("Worker %d: analyzed %d tweets", workerID, len(batchResults))

		results <- batchResults
	}

	log.Printf("Worker %d finished", workerID)
}
