package gemini

import (
	"context"

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

func (c *Client) GenerateText(ctx context.Context, text string) (string, error) {
	config := &genai.GenerateContentConfig{SystemInstruction: genai.NewContentFromText(systemInstruction, genai.RoleModel)}

	resp, err := c.client.Models.GenerateContent(ctx, c.model, genai.Text(text), config)
	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}
