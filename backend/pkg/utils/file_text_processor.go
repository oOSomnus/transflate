package utils

import (
	"regexp"
	"strings"
)

// SplitString splits a string into chunks with each chunk containing up to maxWords words and returns a slice of chunks.
func SplitString(s string, maxWords int) []string {
	words := strings.Fields(s)
	var chunks []string
	for i := 0; i < len(words); i += maxWords {
		end := i + maxWords
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, strings.Join(words[i:end], " "))
	}
	return chunks
}

// GetLastNWords extracts the last n words from the given input string and returns them as a single string.
// If n is greater than or equal to the total number of words, the entire input string is returned.
func GetLastNWords(input string, n int) string {
	words := strings.Fields(input)
	if n >= len(words) {
		return strings.Join(words, " ")
	}
	return strings.Join(words[len(words)-n:], " ")
}

// RemoveNonUnicodeCharacters removes non-Unicode characters from a given string and returns the cleaned string.
func RemoveNonUnicodeCharacters(input string) string {
	re := regexp.MustCompile(`[^\x{0000}-\x{10FFFF}]+`)
	return re.ReplaceAllString(input, "")
}

// ReplaceMultipleSpaces replaces multiple consecutive spaces in the input string with a single space.
func ReplaceMultipleSpaces(input string) string {
	re := regexp.MustCompile(`\s{2,}`)
	return re.ReplaceAllString(input, " ")
}

// TextCleaning processes a string by removing non-Unicode characters and replacing multiple spaces with a single space.
func TextCleaning(str string) string {
	str = RemoveNonUnicodeCharacters(str)
	str = ReplaceMultipleSpaces(str)
	return str
}
