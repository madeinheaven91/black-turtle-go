package errors

import (
	"fmt"

	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
)

type BotError struct {
	Err      error
	Message  string
	ErrorKey lexicon.ErrorKey
	Metadata map[string]any
}

func (e BotError) Error() string {
	return fmt.Sprintf("%s: %v (metadata: %+v)", e.Message, e.Err, e.Metadata)
}

// Create a BotError from error
func From(err error, message string, key lexicon.ErrorKey, metadata map[string]any) *BotError {
	return &BotError{
		Err:      err,
		Message:  message,
		ErrorKey: key,
		Metadata: metadata,
	}
}

// Get display message
func (e *BotError) Display() string {
	return lexicon.Error(e.ErrorKey)
}

// Get display error message from lexicon with key with metadata
func Get(key lexicon.ErrorKey) string {
	return lexicon.Error(key)
}
