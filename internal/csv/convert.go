package csv

import (
	"encoding/csv"
	"strconv"
	"strings"

	"github.com/abdulmuminakinde/tweet-audit/internal/gemini"
)

func ConvertToCSV(results []gemini.TweetAnalysisResult) (string, error) {
	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	headers := []string{"Tweet URL", "Flag", "Action", "Reason"}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	for _, res := range results {
		row := []string{
			res.TweetUrl,
			strconv.FormatBool(res.Deleted),
			res.Action,
			res.Reason,
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	return sb.String(), nil
}
