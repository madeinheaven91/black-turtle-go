package errors

import "fmt"

type DBError struct {
	message string
}

func (e DBError) Error() string {
	return fmt.Sprintf("db error: %q", e.message)
}

func NewDBError(message string) DBError {
	return DBError{message: message}
}
