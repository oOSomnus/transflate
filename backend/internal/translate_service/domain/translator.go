package domain

/*
Translator defines an interface for translation functionality.

Methods:
  - Translate(prevContext, text string) (string, error): Translates the provided text into another language, using the given context for reference.
*/
type Translator interface {
	Translate(prevContext, text string) (string, error)
}
