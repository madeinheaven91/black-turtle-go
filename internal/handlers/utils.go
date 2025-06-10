package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/errors"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
)

func handleBotError(err error, ctx context.Context, b *bot.Bot, update *botmodels.Update) bool {
	if err != nil {
		e, ok := err.(*errors.BotError)
		if ok {
			logging.Error("%q", e)
			b.SendMessage(ctx, params(update, e.Display()))
		} else {
			logging.Error("%q", err)
			b.SendMessage(ctx, params(update, lexicon.Error(lexicon.EGeneral)))
		}
		return false
	}
	return true
}

func params(update *botmodels.Update, text string) *bot.SendMessageParams {
	return &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: text,
		ParseMode: botmodels.ParseModeHTML,
	}

}

func addReplyMarkup(params *bot.SendMessageParams, kb botmodels.ReplyMarkup) *bot.SendMessageParams {
	params.ReplyMarkup = kb
	return params
}
