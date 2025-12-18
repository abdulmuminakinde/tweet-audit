package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	cmd "github.com/abdulmuminakinde/tweet-audit/cmd/ingest"
	"github.com/abdulmuminakinde/tweet-audit/internal/core"
	"github.com/abdulmuminakinde/tweet-audit/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is environment variable is not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	dbQueries := database.New(dbConn)

	cfg := cmd.Config{
		DB:      dbConn,
		Queries: dbQueries,
	}

	file, err := core.LoadTweets("./internal/tweets.js")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = cfg.StreamTweets(ctx, file)
	if err != nil {
		log.Fatal(err)
	}

}
