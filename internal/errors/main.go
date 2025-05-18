package errors

import (
	"fmt"

	"github.com/madeinheaven91/black-turtle-go/internal/logging"
)

type BotError struct {
	message string
}

func (e BotError) Error() string {
	return fmt.Sprintf("bot error: %s", e.message)
}

func NewBotError(message string) BotError {
	return BotError{message: message}
}

func Log(e error, critical ...bool) {
	if len(critical) == 0 {
		logging.Error("%s\n", e)
	} else if len(critical) == 1 && critical[0] {
		logging.Critical("%s\n", e)
	}
}
