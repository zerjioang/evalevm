package parser

import (
	"errors"
	"strings"
)

// ExtractBetween extracts a substring from input between startWord and endWord.
// Returns an error if startWord or endWord are not found, or if endWord comes before startWord.
func ExtractBetween(input, startWord, endWord string) (string, error) {
	startIdx := strings.Index(input, startWord)
	if startIdx == -1 {
		return "", errors.New("start word not found")
	}

	endIdx := strings.Index(input, endWord)
	if endIdx == -1 {
		return "", errors.New("end word not found")
	}

	// Ensure endWord comes after startWord
	if endIdx <= startIdx {
		return "", errors.New("end word appears before start word")
	}

	// Extract substring between them
	startPos := startIdx + len(startWord)
	return input[startPos:endIdx], nil
}
