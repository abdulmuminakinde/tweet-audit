package gemini

import (
	"context"

	"google.golang.org/genai"
)

func (c *Client) GenerateText(ctx context.Context) (string, error) {
	resp, err := c.client.Models.GenerateContent(ctx, c.model, genai.Text("What should a good golang project structure look like?"), nil)
	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}
