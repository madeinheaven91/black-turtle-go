package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/internal/messages"
	"github.com/madeinheaven91/black-turtle-go/internal/query/ir"
	"github.com/madeinheaven91/black-turtle-go/internal/query/parser"
	"github.com/madeinheaven91/black-turtle-go/internal/requests"
	"github.com/madeinheaven91/black-turtle-go/pkg/errors"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func LessonsMatch(update *models.Update) bool {
	return strings.HasPrefix(strings.ToLower(update.Message.Text), "пары")
}

func LessonsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// 1. parse input into request struct
	// 2. fetch week based on request
	// 3. select requested day from a week
	// 4. build a message
	// 5. send a message

	logging.Debug("Handling lesson request")
	input := update.Message.Text
	p := parser.FromString(input)
	raw := p.ParseQuery()
	if len(p.Errors()) != 0 {
		logging.Error(fmt.Sprintf("parser errors: %q\n", p.Errors()))
		b.SendMessage(ctx, params(update, errors.Get(lexicon.EParser)))
		return
	}
	lqr, ok := (*raw).(*ir.LessonsQueryRaw)
	if !ok {
		logging.Error("can't convert *QueryRaw into *LessonsQueryRaw, wtf??")
		return
	}
	query, err := lqr.Validate(update.Message.Chat.ID)
	ok = handleBotError(err, ctx, b, update)
	if !ok {
		return
	}

	resp, err := requests.FetchWeek(query)
	ok = handleBotError(err, ctx, b, update)
	if !ok {
		return
	}

	week := resp.IntoWeek()

	var displayName string
	if resp.Group != nil {
		displayName = resp.Group.Name
	} else if resp.Teacher != nil {
		displayName = *resp.Teacher.Fio
	} else {
		displayName = "null"
	}
	msg := messages.BuildDayMsg(
		week.Days[shared.NormalizeWeekday(query.Date.Weekday())],
		query.Date,
		displayName,
	)

	b.SendMessage(ctx, params(update, msg))
	logging.Trace("Done handling lesson request")
}
