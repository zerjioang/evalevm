package parser_test

import (
	"evalevm/internal/parser"
	"testing"
)

func TestExtractBetween(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		startWord string
		endWord   string
		want      string
		wantErr   bool
	}{
		{
			name:      "happy path",
			input:     "Hello this is the start of something and here is the end marker.",
			startWord: "start",
			endWord:   "end",
			want:      " of something and here is the ",
			wantErr:   false,
		},
		{
			name:      "start word missing",
			input:     "Hello world with only end marker",
			startWord: "start",
			endWord:   "end",
			want:      "",
			wantErr:   true,
		},
		{
			name:      "end word missing",
			input:     "Hello world with only start marker",
			startWord: "start",
			endWord:   "end",
			want:      "",
			wantErr:   true,
		},
		{
			name:      "end before start",
			input:     "end comes before start in this string",
			startWord: "start",
			endWord:   "end",
			want:      "",
			wantErr:   true,
		},
		{
			name:      "adjacent words",
			input:     "startend",
			startWord: "start",
			endWord:   "end",
			want:      "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ExtractBetween(tt.input, tt.startWord, tt.endWord)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractBetween() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractBetween() got = %q, want %q", got, tt.want)
			}
		})
	}
}
