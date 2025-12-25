package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/abdulmuminakinde/tweet-audit/internal/database"
	"google.golang.org/genai"
)

const systemInstruction = `You are a tweet content evaluator for professional social media auditing.

Your task is to analyze tweets and determine if they should be flagged for deletion based on alignment criteria.

Evaluate each tweet against these categories:
1. Professionalism: Language, tone, word choice
2. Keywords: Presence of specific flagged terms
3. Outdated Content: References to old/changed opinions, obsolete information
4. Brand Alignment: Consistency with current professional image

For each tweet, respond in this exact format:
ACTION: [KEEP/FLAG]
CATEGORY: [professionalism/keywords/outdated/brand]
SEVERITY: [low/medium/high]
REASON: [One sentence explanation]

Be strict but fair. When in doubt, flag for review.
Focus on content that could harm professional reputation.`

type TweetAnalysisResult struct {
	TweetUrl string `json:"tweet_url"`
	Action   string `json:"action"`
	Deleted  bool   `json:"deleted"`
	Reason   string `json:"reason"`
}

func (c *Client) GenerateText(ctx context.Context, text string) (string, error) {
	config := &genai.GenerateContentConfig{SystemInstruction: genai.NewContentFromText(systemInstruction, genai.RoleModel)}

	resp, err := c.client.Models.GenerateContent(ctx, c.model, genai.Text(text), config)
	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}

func (c *Client) AnalyzeTweets(ctx context.Context, tweets []database.GetTweetsRow) ([]TweetAnalysisResult, error) {
	prompt := buildPrompt(tweets)

	part := genai.NewPartFromText(systemInstruction)

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"tweet_url": {
						Type:        genai.TypeString,
						Description: "The URL of the tweet",
					},
					"action": {
						Type:        genai.TypeString,
						Description: "KEEP or FLAG for deletion",
					},
					"deleted": {
						Type:        genai.TypeBoolean,
						Description: "Tweet deleted or not deleted",
					},
					"reason": {
						Type:        genai.TypeString,
						Description: "Brief explanation for the decision",
					},
				},
				Required: []string{"tweet_url", "action", "reason"},
			},
		},
		SystemInstruction: &genai.Content{Parts: []*genai.Part{part}},
	}
	resp, err := c.client.Models.GenerateContent(ctx, c.model, genai.Text(prompt), config)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze batch, %w", err)
	}

	var results []TweetAnalysisResult
	if err := json.Unmarshal([]byte(resp.Text()), &results); err != nil {
		return nil, fmt.Errorf("failed to parse the response: %w", err)
	}

	return results, nil
}

func buildPrompt(tweets []database.GetTweetsRow) string {
	var sb strings.Builder

	sb.WriteString("Analyze these tweets for potential deletion. ")
	sb.WriteString("Evaluate each against professionalism, brand alignment, and content quality.\n\n")

	for i, tweet := range tweets {
		sb.WriteString(fmt.Sprintf("Tweet %d:\n", i+1))
		sb.WriteString(fmt.Sprintf("Text: %s\n", tweet.FullText))
		sb.WriteString(fmt.Sprintf("URL: %s\n", tweet.Url))
		sb.WriteString(fmt.Sprintf("Retweeted: %v", tweet.Retweeted))
	}

	sb.WriteString("For each tweet, determine if it should be KEPT or FLAGGED for deletion.")

	return sb.String()
}
