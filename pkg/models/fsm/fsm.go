package fsm

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type (
	State        int
	StateHandler func(ctx context.Context, b *bot.Bot, update *models.Update)
)

const (
	EnterGroup State = iota
	EnterTeacher
	RegCancel
)
