package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/abdulmuminakinde/tweet-audit/internal/config"
	"github.com/abdulmuminakinde/tweet-audit/internal/gemini"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter tweet to be analyzed: ")

	question, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	trimmedQuestion := strings.TrimSpace(question)

	ctx := context.Background()
	cfg, err := config.LoadOrCreateConfig()
	if err != nil {
		log.Fatal("error loading config")
	}

	geminiclient, err := gemini.NewClient(ctx, cfg)
	if err != nil {
		log.Fatal("error creating gemini client")
	}

	result, err := geminiclient.GenerateText(ctx, trimmedQuestion)
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}

	fmt.Println(result)
}
