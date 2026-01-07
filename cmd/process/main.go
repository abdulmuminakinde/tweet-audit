package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/abdulmuminakinde/tweet-audit/internal/config"
	"github.com/abdulmuminakinde/tweet-audit/internal/database"
	"github.com/abdulmuminakinde/tweet-audit/internal/gemini"

	_ "github.com/lib/pq"
)

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	cfg, err := config.LoadOrCreateConfig()
	if err != nil {
		log.Fatal("error loading config")
	}

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

	processor := NewTweetProcessor(client, dbQueries)

	processor.ProcessAllTweets(ctx)
}
