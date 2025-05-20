package errors

import (
	"fmt"
)

type BotError struct {
	Err        error
	Message    string
	Metadata   map[string]any
}

func (e BotError) Error() string {
	return fmt.Sprintf("%s: %v (metadata: %+v)", e.Message, e.Err, e.Metadata)
}

func Wrap(err error, message string, metadata map[string]any) *BotError {
	return &BotError {
		Err: err,
		Message: message,
		Metadata: metadata,
	}
}
