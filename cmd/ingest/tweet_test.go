package cmd

import (
	"reflect"
	"testing"
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
