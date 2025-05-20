package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

func reply(ctx context.Context, b *bot.Bot, update *botmodels.Update, message string) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      message,
		ParseMode: botmodels.ParseModeHTML,
	})
}

