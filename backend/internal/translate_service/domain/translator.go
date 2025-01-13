package domain

// Translator is an interface for handling text translation with context awareness.
// Translate translates the provided text based on the previous context and returns the result or an error if any.
type Translator interface {
	Translate(prevContext, text string) (string, error)
}
