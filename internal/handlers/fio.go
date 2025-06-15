package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/internal/messages"
	"github.com/madeinheaven91/black-turtle-go/internal/query/ir"
	"github.com/madeinheaven91/black-turtle-go/internal/query/parser"
	"github.com/madeinheaven91/black-turtle-go/pkg/errors"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func FioMatch(update *botmodels.Update) bool {
	if update.Message == nil {
		return false
	} else {
		return strings.HasPrefix(strings.ToLower(update.Message.Text), "фио")
	}
}

func FioHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	logging.Debug("Handling FIO request")
	p := parser.FromString(update.Message.Text)
	rawQ := p.ParseQuery()
	if len(p.Errors()) != 0 {
		logging.Error(fmt.Sprintf("parser errors: %q\n", p.Errors()))
		b.SendMessage(ctx, shared.Params(update, errors.Get(lexicon.EParser)))
		return
	}
	rawQuery, ok := (*rawQ).(*ir.FioQueryRaw)
	if !ok {
		logging.Critical("can't convert *QueryRaw into *FioQueryRaw, wtf??")
		return
	}
	query := rawQuery.Validate()
	entities, err := db.GetStudyEntities(query.Name)
	if err != nil {
		logging.Error("%s", err.Error())
		b.SendMessage(ctx, shared.Params(update, lexicon.Error(lexicon.EGeneral)))
	}
	b.SendMessage(ctx, shared.Params(update, messages.BuildFIOMessage(entities)))
	logging.Trace("Done handling FIO request")
}
