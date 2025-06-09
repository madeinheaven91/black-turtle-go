package handlers

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/keyboards"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
)

func HelpMatch(update *botmodels.Update) bool {
	return strings.HasPrefix(strings.ToLower(update.Message.Text), "помощь")
}

func HelpHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      lexicon.Get(lexicon.HelpGeneric),
		ReplyMarkup: keyboards.Help(),
		ParseMode: botmodels.ParseModeHTML,
	})
}
