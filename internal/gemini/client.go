package gemini

import (
	"context"
	"fmt"

	"github.com/abdulmuminakinde/tweet-audit/internal/config"
	"google.golang.org/genai"
)

type Client struct {
	client *genai.Client
	model  string
}

func NewClient(ctx context.Context, config *config.Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}

	if config.AIModel == "" {
		config.AIModel = "gemini-2.5-flash"
	}

	genaiclient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client: genaiclient,
		model:  config.AIModel,
	}, nil
}
