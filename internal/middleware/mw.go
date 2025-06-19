package middleware

import (
	"context"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"

	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func LogRequest(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		name := shared.GetFromName(update)
		if update.Message.Chat.Type == "private" {
			logging.Telegram("%s (%d): %s\n", name, update.Message.Chat.ID, update.Message.Text)
		} else {
			chat_name := shared.GetChatName(update)
			logging.Telegram("%s (%d) in %s (%d): %s\n", name, update.Message.From.ID, chat_name, update.Message.Chat.ID, update.Message.Text)
		}
		next(ctx, b, update)
	}
}

func DbSync(next bot.HandlerFunc) bot.HandlerFunc {
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
