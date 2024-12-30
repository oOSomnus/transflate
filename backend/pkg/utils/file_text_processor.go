package utils

import (
	"regexp"
	"strings"
)

/*
SplitString splits a given string into chunks containing up to a specified maximum number of words.

Parameters:
  - s (string): The input string to be split.
  - maxWords (int): The maximum number of words allowed in each chunk.

Returns:
  - ([]string): A slice of strings, where each string contains up to `maxWords` words from the input string.
*/

func SplitString(s string, maxWords int) []string {
	words := strings.Fields(s)
	var chunks []string
	for i := 0; i < len(words); i++ {
		end := i + maxWords
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, strings.Join(words[i:end], " "))
	}
	return chunks
}

/*
GetLastNWords extracts the last N words from the given input string.

Parameters:
  - input (string): The input string from which words are extracted.
  - n (int): The number of words to extract from the end of the input string.

Returns:
  - (string): A string containing the last N words, or all words if N exceeds the total word count.
*/

func GetLastNWords(input string, n int) string {
	words := strings.Fields(input)
	if n >= len(words) {
		return strings.Join(words, " ")
	}
	return strings.Join(words[len(words)-n:], " ")
}

func RemoveNonUnicodeCharacters(input string) string {
	re := regexp.MustCompile(`[^\x{0000}-\x{10FFFF}]+`)
	return re.ReplaceAllString(input, "")
}

func ReplaceMultipleSpaces(input string) string {
	re := regexp.MustCompile(`\s{2,}`)
	return re.ReplaceAllString(input, " ")
}
