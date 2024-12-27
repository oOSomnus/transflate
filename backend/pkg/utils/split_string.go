package utils

import "strings"

/*
SplitString splits a given string into smaller chunks, ensuring each chunk does not exceed the specified maximum token length.

Parameters:
  - s (string): The input string to be split into chunks.
  - maxTokens (int): The maximum number of tokens (characters including spaces) allowed per chunk.

Returns:
  - ([]string): A slice of strings where each element is a chunk of the original string, adhering to the token limit.
*/
func SplitString(s string, maxTokens int) []string {
	words := strings.Fields(s)
	var chunks []string
	var chunk []string
	var count int

	for _, word := range words {
		if count+len(word)+1 > maxTokens {
			chunks = append(chunks, strings.Join(chunk, " "))
			chunk = nil
			count = 0
		}
		chunk = append(chunk, word)
		count += len(word) + 1
	}

	if len(chunk) > 0 {
		chunks = append(chunks, strings.Join(chunk, " "))
	}

	return chunks
}
