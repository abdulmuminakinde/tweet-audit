package core

import (
	"os"
)

type Post struct {
	Tweet struct {
		Retweeted         bool   `json:"retweeted"`
		Truncated         bool   `json:"truncated"`
		ID                string `json:"id"`
		PossiblySensitive bool   `json:"possibly_sensitive"`
		CreatedAt         string `json:"created_at"`
		FullText          string `json:"full_text"`
	} `json:"tweet"`
}

func LoadTweets(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
