package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/internal/query/ir"
	"github.com/madeinheaven91/black-turtle-go/internal/query/parser"
	"github.com/madeinheaven91/black-turtle-go/pkg/errors"
	"github.com/madeinheaven91/black-turtle-go/pkg/keyboards"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func HelpMatch(update *botmodels.Update) bool {
	if update.Message == nil {
		return false
	} else {
		return strings.HasPrefix(strings.ToLower(update.Message.Text), "помощь")
	}
}

func (c Container) HelpHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	logging.Debug("Handling help request")
	p := parser.FromString(update.Message.Text)
	rawQ := p.ParseQuery()
	if len(p.Errors()) != 0 {
		logging.Error(fmt.Sprintf("parser errors: %q\n", p.Errors()))
		b.SendMessage(ctx, shared.Params(update, errors.Get(lexicon.EParser)))
		return
	}
	rawQuery, ok := (*rawQ).(*ir.HelpQueryRaw)
	if !ok {
		logging.Critical("can't convert *QueryRaw into *HelpQueryRaw, wtf??")
		return
	}
	query := rawQuery.Validate()
	switch query.Command {
	case models.Nil:
		b.SendMessage(ctx, shared.AddReplyMarkup(shared.Params(update, lexicon.Get(lexicon.HelpGeneric)), keyboards.Help()))
	case models.Lessons:
		b.SendMessage(ctx, shared.AddReplyMarkup(shared.Params(update, lexicon.Get(lexicon.HelpLessons)), keyboards.Help()))
	case models.Bells:
		b.SendMessage(ctx, shared.AddReplyMarkup(shared.Params(update, lexicon.Get(lexicon.HelpBells)), keyboards.Help()))
	case models.Fio:
		b.SendMessage(ctx, shared.AddReplyMarkup(shared.Params(update, lexicon.Get(lexicon.HelpFio)), keyboards.Help()))
	}

	logging.Trace("Done handling help request")
}
