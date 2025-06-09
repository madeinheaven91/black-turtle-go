package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/pkg/keyboards"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
)

func StartHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	conn := db.GetConnection()
	defer db.CloseConn(conn)
	err := db.AddChat(conn, update)
	if err != nil {
		logging.Error("%s\n", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        lexicon.Error(lexicon.EGeneral),
			ParseMode:   botmodels.ParseModeHTML,
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "Добавил",
			ReplyMarkup: keyboards.Default(),
			ParseMode: botmodels.ParseModeHTML,
		})
	}
}
