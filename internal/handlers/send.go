package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

func (c Container) SendHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	c.Notif.SendHandler(ctx, b, update)
}
