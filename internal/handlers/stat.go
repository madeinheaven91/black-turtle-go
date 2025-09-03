package handlers

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func (c *Container) StatHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	// TODO: full fetched monitoring engine with more metrics and better reports
	private, groups, err := db.ChatCount()
	if err != nil {
		b.SendMessage(ctx, shared.Params(update, lexicon.Error(lexicon.EGeneral)))
	} else {
		msg := strings.Join([]string{
			"<b>Статистика</b>",
			"-----------",
			"Всего чатов: " + strconv.Itoa(private+groups),
			"Личных: " + strconv.Itoa(private),
			"Групп: " + strconv.Itoa(groups),
		}, "\n")
		b.SendMessage(ctx, shared.Params(update, msg))
	}
}
