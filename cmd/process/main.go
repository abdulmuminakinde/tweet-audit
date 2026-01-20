package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	"github.com/abdulmuminakinde/tweet-audit/internal/config"
	"github.com/abdulmuminakinde/tweet-audit/internal/database"
	"github.com/abdulmuminakinde/tweet-audit/internal/gemini"

	_ "github.com/lib/pq"
)

func main() {

	setBatchSize := flag.Int("setbatchsize", 400, "The batch size for tweets")
	setNumWorkers := flag.Int("setnumworkers", 3, "Number of workers to analyze tweets")
	setLimiter := flag.Int("limiter", 12, "The number of seconds between API calls")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	cfg, err := config.LoadOrCreateConfig()

	if err != nil {
		log.Fatal("error loading config")
	}

	cfg.BatchSize = *setBatchSize
	cfg.NumWorkers = *setNumWorkers
	cfg.Limiter = *setLimiter

	client, err := gemini.NewClient(ctx, cfg)
	if err != nil {
		log.Fatal("error creating gemini client")
	}

	dbURL := os.Getenv("DATABASE_URL")

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	dbQueries := database.New(dbConn)

	processor := NewTweetProcessor(client, dbQueries, cfg)

	processor.ProcessAllTweets(ctx)
}
