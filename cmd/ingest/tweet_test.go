package cmd

import (
	"reflect"
	"testing"
	"time"

	"github.com/abdulmuminakinde/tweet-audit/internal/core"
)

func TestNormalizeText(t *testing.T) {
	tests := map[string]struct {
		text   string
		output string
	}{
		"with trailing url":    {text: "This is a tweet with a trailing url https://t.co/dkjdh", output: "This is a tweet with a trailing url"},
		"without trailing url": {text: "This is a tweet without a trailing url", output: "This is a tweet without a trailing url"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := normalizeText(tc.text)
			if !reflect.DeepEqual(tc.output, got) {
				t.Fatalf("expected: %v, got: %v", tc.text, got)
			}
		})
	}
}

func TestGetTweetUrl(t *testing.T) {
	tests := map[string]struct {
		tweetID   string
		want      string
		username  string
		wantError bool
	}{
		"with a tweet ID": {
			tweetID:   "1234567890",
			want:      "https://x.com/lanrey_waju/status/1234567890",
			username:  "lanrey_waju",
			wantError: false,
		},
		"with an empty tweet ID": {
			tweetID:   "",
			want:      "",
			wantError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tweet := core.RawTweet{}
			tweet.Tweet.ID = tc.tweetID

			got, err := getTweetUrl(tweet, tc.username)

			if tc.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %q, got: %q", tc.want, got)
			}
		})
	}
}

func TestParseTweetDate(t *testing.T) {
	tests := map[string]struct {
		want      time.Time
		wantError bool
		createdAt string
	}{
		"valid tweet date": {
			createdAt: "Sat Dec 20 02:42:20 +0000 2025",
			want:      time.Date(2025, 12, 20, 02, 42, 20, 0, time.UTC),
		},
		"empty date string": {
			createdAt: "",
			wantError: true,
		},
		"nonsensical input": {
			createdAt: "anything but date",
			wantError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := parseTweetDate(tc.createdAt)

			if tc.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !got.Equal(tc.want) {
				t.Errorf("want: %v, got: %v", tc.want, got)
			}
		})
	}

}
