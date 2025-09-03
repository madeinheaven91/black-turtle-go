package middleware

import (
	"context"
	"runtime/debug"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"

	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/pkg/errors"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func LogRequest(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		name := shared.GetFromName(update)
		if update.Message.Chat.Type == "private" {
			logging.Telegram("%s (%d): %s\n", name, update.Message.Chat.ID, update.Message.Text)
		} else {
			chatName := shared.GetChatName(update)
			logging.Telegram("%s (%d) in %s (%d): %s\n", name, update.Message.From.ID, chatName, update.Message.Chat.ID, update.Message.Text)
		}
		next(ctx, b, update)
	}
}

func DBSync(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		err := db.SyncName(update)
		if err != nil {
			logging.Error("%s\n", err)
			next(ctx, b, update)
		} else {
			next(ctx, b, update)
		}
	}
}

func Recover(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logging.Critical("panic: %s\n%s\n", rvr, debug.Stack())
				b.SendMessage(ctx, shared.Params(update, errors.Get(lexicon.EParser)))
			}
		}()
		next(ctx, b, update)
	}
}

func CheckAdmin(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		chatID := shared.GetChatID(update)
		isAdmin := db.CheckAdmin(chatID)
		if isAdmin {
			next(ctx, b, update)
		} else {
			b.SendMessage(ctx, shared.Params(update, "Команда доступна только админам"))
		}
	}
}
