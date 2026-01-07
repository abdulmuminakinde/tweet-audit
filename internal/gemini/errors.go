package gemini

import (
	"fmt"
	"strings"
)

func (c *Client) shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	errorString := strings.ToLower(err.Error())

	if strings.Contains(errorString, "invalid api key") ||
		strings.Contains(errorString, "unauthorized") ||
		strings.Contains(errorString, "context canceled") {
		return false
	}

	if strings.Contains(errorString, "rate limit") ||
		strings.Contains(errorString, "timeout") ||
		strings.Contains(errorString, "temporary") ||
		strings.Contains(errorString, "503") ||
		strings.Contains(errorString, "429") {
		return true
	}

	// retry unknown errors
	return true
}

func (c *Client) wrapError(err error) error {
	if err == nil {
		return nil
	}

	errorString := strings.ToLower(err.Error())

	switch {
	case strings.Contains(errorString, "rate limit") || strings.Contains(errorString, "429"):
		return fmt.Errorf("rate limit exceeded (5 RPM): %w", err)
	case strings.Contains(errorString, "quota"):
		return fmt.Errorf("daily quota exceeded (20 RPD): %w", err)
	case strings.Contains(errorString, "invalid api key"):
		return fmt.Errorf("invalid API key (check GEMINI_API_KEY): %w", err)
	case strings.Contains(errorString, "timeout"):
		return fmt.Errorf("request timeout (batch too large?): %w", err)
	default:
		return err
	}
}
