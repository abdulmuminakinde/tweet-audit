package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"

	cmd "github.com/abdulmuminakinde/tweet-audit/cmd/ingest"
	"github.com/abdulmuminakinde/tweet-audit/internal/config"
	"github.com/abdulmuminakinde/tweet-audit/internal/core"
	"github.com/abdulmuminakinde/tweet-audit/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	setUsername := flag.String("setusername", "", "The X username to construct tweet url")
	setApiKey := flag.String("setapikey", "", "Google AI Studio API Key")

	flag.Parse()

	if *setUsername != "" || *setApiKey != "" {
		cfg, err := config.LoadOrCreateConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		if *setUsername != "" {
			cfg.Username = *setUsername
			log.Println("Username updated")
		}
		if *setApiKey != "" {
			cfg.APIKey = *setApiKey
			log.Println("API Key updated")
		}

		if err := cfg.Save(); err != nil {
			log.Fatalf("error saving config: %v", err)
		}

		log.Println("Config saved successfully")
	}

	config, err := config.LoadOrCreateConfig()
	if err != nil {
		log.Fatal("error reading config")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	dbQueries := database.New(dbConn)

	s := cmd.State{
		Config:  config,
		DB:      dbConn,
		Queries: dbQueries,
	}

	file, err := core.LoadTweets("./internal/tweets.js")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = s.StreamTweets(ctx, file)
	if err != nil {
		log.Fatal(err)
	}

}
