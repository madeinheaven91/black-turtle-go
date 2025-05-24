package errors

import (
	"fmt"

	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
)

type BotError struct {
	Err        error
	Message    string
	LexiconKey string
	Metadata   map[string]any
}

func (e BotError) Error() string {
	return fmt.Sprintf("%s: %v (metadata: %+v)", e.Message, e.Err, e.Metadata)
}

// Create a BotError from error
func From(err error, message string, key string, metadata map[string]any) *BotError {
	return &BotError{
		Err:        err,
		Message:    message,
		LexiconKey: key,
		Metadata:   metadata,
	}
}

// Get display message
func (e *BotError) Display() string {
	return lexicon.Error(e.LexiconKey)
}

// Get general display error message
func General() string {
	return lexicon.Error("")
}

// Get display error message from lexicon with key with metadata
func Get(key string, metadata ...any) string {
	return lexicon.Error(key, metadata...)
}
