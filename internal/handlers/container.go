package handlers

import (
	"github.com/go-telegram/bot"
	"github.com/madeinheaven91/black-turtle-go/internal/handlers/admin"
)

type Container struct {
	Bot     *bot.Bot
	Notif *admin.NotificationService
}

func NewContainer(b *bot.Bot, notif *admin.NotificationService) *Container {
	return &Container{
		Bot:     b,
		Notif: notif,
	}
}
