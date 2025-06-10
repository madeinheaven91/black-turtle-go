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
		b.SendMessage(ctx, params(update, lexicon.Error(lexicon.EGeneral)))
	} else {
		b.SendMessage(ctx, addReplyMarkup(params(update, "Добавил"), keyboards.Default()))
	}
}
