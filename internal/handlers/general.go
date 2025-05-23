package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/internal/errors"
	"github.com/madeinheaven91/black-turtle-go/internal/logging"
)

func reply(message string, ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      message,
		ParseMode: botmodels.ParseModeHTML,
	})
}


func handleBotError(err error, ctx context.Context, b *bot.Bot, update *botmodels.Update) bool {
	if err != nil {
		e, ok := err.(*errors.BotError)
		if ok{
			logging.Error("%q", e)
			reply(e.Display(), ctx, b, update)
		} else {
			logging.Error("%q", err)
			reply(errors.General(), ctx, b, update)
		}
		return false
	}
	return true
}
