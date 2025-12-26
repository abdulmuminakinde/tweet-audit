package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/abdulmuminakinde/tweet-audit/internal/config"
	"github.com/abdulmuminakinde/tweet-audit/internal/database"
	"github.com/abdulmuminakinde/tweet-audit/internal/gemini"

	_ "github.com/lib/pq"
)

func main() {

	ctx := context.Background()
	cfg, err := config.LoadOrCreateConfig()
	if err != nil {
		log.Fatal("error loading config")
	}

	geminiclient, err := gemini.NewClient(ctx, cfg)
	if err != nil {
		log.Fatal("error creating gemini client")
	}

	dbURL := os.Getenv("DATABASE_URL")

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	dbQueries := database.New(dbConn)

	tweets, err := dbQueries.GetTweets(ctx)
	if err != nil {
		log.Fatal(err)
	}

	gemini.ProcessBatchedTweets(ctx, geminiclient, tweets)
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}

}
