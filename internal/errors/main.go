package errors

import (
	"fmt"

	"github.com/madeinheaven91/black-turtle-go/internal/lexicon"
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

func From(err error, message string, key string, metadata map[string]any) *BotError {
	return &BotError{
		Err:        err,
		Message:    message,
		LexiconKey: key,
		Metadata:   metadata,
	}
}

func (e *BotError) Display() string {
	return lexicon.Error(e.LexiconKey)
}

func General() string {
	return lexicon.Error("")
}
