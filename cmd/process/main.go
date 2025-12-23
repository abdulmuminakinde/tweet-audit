package main

import (
	"context"
	"fmt"
	"log"

	"github.com/abdulmuminakinde/tweet-audit/internal/config"
	"github.com/abdulmuminakinde/tweet-audit/internal/gemini"
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

	result, err := geminiclient.GenerateText(ctx)
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}

	fmt.Println(result)
}
