package utils

import (
	"reflect"
	"testing"
)

func TestSplitString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxWords int
		expected []string
	}{
		{"short input", "hello world", 5, []string{"hello world"}},
		{"exact words", "one two three four five", 5, []string{"one two three four five"}},
		{"split evenly", "one two three four five six", 3, []string{"one two three", "four five six"}},
		{"split remainder", "one two three four five", 2, []string{"one two", "three four", "five"}},
		{"empty input", "", 2, []string{""}},
		{"single chunk", "one", 2, []string{"one"}},
	}

	for _, tc := range tests {
		t.Run(
			tc.name, func(t *testing.T) {
				result := SplitString(tc.input, tc.maxWords)
				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			},
		)
	}
}

func TestGetLastNWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		n        int
		expected string
	}{
		{"all words", "hello world", 5, "hello world"},
		{"exact last n words", "one two three four five", 3, "three four five"},
		{"n is zero", "one two three", 0, ""},
		{"n greater than input length", "hello world", 10, "hello world"},
		{"empty input", "", 3, ""},
		{"last word", "one", 1, "one"},
	}

	for _, tc := range tests {
		t.Run(
			tc.name, func(t *testing.T) {
				result := GetLastNWords(tc.input, tc.n)
				if result != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, result)
				}
			},
		)
	}
}

func TestRemoveNonUnicodeCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"valid unicode", "こんにちは", "こんにちは"},
		{"mixed characters", "hello世界\x83", "hello世界"},
		{"non-unicode chars", "\x80\x81\x82", ""},
		{"empty input", "", ""},
		{"numeric input", "12345", "12345"},
	}

	for _, tc := range tests {
		t.Run(
			tc.name, func(t *testing.T) {
				result := RemoveNonUnicodeCharacters(tc.input)
				if result != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, result)
				}
			},
		)
	}
}

func TestReplaceMultipleSpaces(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no spaces", "test", "test"},
		{"single spaces", "hello world", "hello world"},
		{"multiple spaces", "hello   world", "hello world"},
		{"leading trailing spaces", "   hello   world   ", " hello world "},
		{"empty input", "", ""},
	}

	for _, tc := range tests {
		t.Run(
			tc.name, func(t *testing.T) {
				result := ReplaceMultipleSpaces(tc.input)
				if result != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, result)
				}
			},
		)
	}
}

func TestTextCleaning(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic cleanup", "hello   world!!  ", "hello world!! "},
		{"unicode and spaces", "こんにちは   世界\x91", "こんにちは 世界"},
		{"only spaces", "     ", " "},
		{"empty input", "", ""},
		{"complex input", "\x80 hello   世界 \x81", " hello 世界 "},
	}

	for _, tc := range tests {
		t.Run(
			tc.name, func(t *testing.T) {
				result := TextCleaning(tc.input)
				if result != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, result)
				}
			},
		)
	}
}
