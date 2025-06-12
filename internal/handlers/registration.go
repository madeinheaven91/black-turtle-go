package handlers

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/keyboards"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
	// "github.com/madeinheaven91/black-turtle-go/internal/db"
	// "github.com/madeinheaven91/black-turtle-go/pkg/keyboards"
	// "github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	// "github.com/madeinheaven91/black-turtle-go/pkg/logging"
)

func RegistrationMatch(update *botmodels.Update) bool {
	if update.Message == nil {
		return false
	} else {
		return strings.HasPrefix(strings.ToLower(update.Message.Text), "регистрация")
	}
}

func Registration(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	b.SendMessage(ctx,
		shared.AddReplyMarkup(
			shared.Params(update, lexicon.Get(lexicon.RegGroupOrTeacher)),
			keyboards.ChooseGroupOrTeacher(),
		),
	)
}

