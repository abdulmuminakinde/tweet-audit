package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/abdulmuminakinde/tweet-audit/internal/gemini"
)

func (p *TweetProcessor) saveResults(ctx context.Context, results []gemini.TweetAnalysisResult) error {
	return p.writeResultsToFile(results)
}

func (p *TweetProcessor) writeResultsToFile(results []gemini.TweetAnalysisResult) error {
	var flagged []gemini.TweetAnalysisResult
	for _, r := range results {
		if r.Action == "FLAG" {
			flagged = append(flagged, r)
		}
	}

	data, err := json.MarshalIndent(flagged, "", "  ")
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("flagged_tweets_%s.json", time.Now().Format("2006-01-02"))
	return os.WriteFile(filename, data, 0644)
}
