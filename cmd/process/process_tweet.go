package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/abdulmuminakinde/tweet-audit/internal/config"
	"github.com/abdulmuminakinde/tweet-audit/internal/csv"
	"github.com/abdulmuminakinde/tweet-audit/internal/database"
	"github.com/abdulmuminakinde/tweet-audit/internal/gemini"
)

type TweetProcessor struct {
	client    *gemini.Client
	dbQueries *database.Queries
	config    *config.Config
}

func NewTweetProcessor(client *gemini.Client, queries *database.Queries, config *config.Config) *TweetProcessor {
	if config.BatchSize == 0 {
		config.BatchSize = 200
	}
	if config.Limiter == 0 {
		config.Limiter = 12
	}
	return &TweetProcessor{
		client:    client,
		dbQueries: queries,
		config:    config,
	}
}

func (p *TweetProcessor) ProcessAllTweets(ctx context.Context) error {
	tweets, err := p.dbQueries.GetTweets(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch tweets: %w", err)
	}

	batchSize := p.config.BatchSize

	batches := chunkTweets(tweets, batchSize)

	results := p.ProcessBatchesConcurrently(ctx, batches)

	log.Printf("Saving %d results...", len(results))
	return p.saveResults(ctx, results)

}

func (p *TweetProcessor) saveCSV(results []gemini.TweetAnalysisResult) error {
	csvData, err := csv.ConvertToCSV(results)
	if err != nil {
		return fmt.Errorf("failed to convert to CSV: %w", err)
	}

	filename := fmt.Sprintf("results_%s.csv", time.Now().Format("2006-01-02"))
	if err := os.WriteFile(filename, []byte(csvData), 0644); err != nil {
		return fmt.Errorf("failed to write CSV: %w", err)
	}

	log.Printf("✓ CSV saved to %s", filename)
	return nil
}

func (p *TweetProcessor) saveResults(ctx context.Context, results []gemini.TweetAnalysisResult) error {
	var wg sync.WaitGroup
	errors := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errors <- p.saveJSON(results)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errors <- p.saveCSV(results)
	}()

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *TweetProcessor) ProcessBatchesConcurrently(ctx context.Context, batches [][]database.GetTweetsRow) []gemini.TweetAnalysisResult {
	const numWorkers = 3

	jobs := make(chan []database.GetTweetsRow, len(batches))
	results := make(chan []gemini.TweetAnalysisResult, len(batches))

	limiter := time.NewTicker(time.Duration(p.config.Limiter) * time.Second) // to match the free tier limit of 5 RPM

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
