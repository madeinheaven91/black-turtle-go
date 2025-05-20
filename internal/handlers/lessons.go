package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/internal/logging"
	"github.com/madeinheaven91/black-turtle-go/internal/messages"
	"github.com/madeinheaven91/black-turtle-go/internal/query/ir"
	"github.com/madeinheaven91/black-turtle-go/internal/query/jsonbuilder"
	"github.com/madeinheaven91/black-turtle-go/internal/query/parser"
	"github.com/madeinheaven91/black-turtle-go/internal/requests"
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
	query := p.ParseQuery()
	if len(p.Errors()) != 0 {
		errMsg := fmt.Sprintf("parser error: %q\n", p.Errors())
		logging.Error(errMsg)
		reply(ctx, b, update, errMsg)
		return
	}
	req, ok := query.Command.(*ir.LessonsQuery)
	if !ok {
		// NOTE: this shouldn't happen, but just in case
		logging.Error("query type assertion error: got %T\n", query)
		reply(ctx, b, update, "Type assertion error")
		return
	}
	json, err := jsonbuilder.BuildPayload(*query, update.Message.Chat.ID)
	if err != nil {
		logging.Error("%q", err)
		reply(ctx, b, update, err.Error())
		return
	}

	// reply(ctx, b, update, json)
	resp, err := requests.FetchWeek(json)
	if err != nil {
		logging.Error("%s\n", err)
		reply(ctx, b, update, err.Error())
		return
	}

	// req, err := parser.ParseToRequest(update.Message.Text)
	// if err != nil {
	// 	logging.Error("%s\n", err)
	// 	b.SendMessage(ctx, &bot.SendMessageParams{
	// 		ChatID:    update.Message.Chat.ID,
	// 		Text:      lexicon.Get("requestError"),
	// 		ParseMode: models.ParseModeHTML,
	// 	})
	// 	return
	// }
	//
	// resp, err := requests.FetchWeek(req.StudyEntityType, req.StudyEntityId, req.Date)
	// if err != nil {
	// 	logging.Error("%s\n", err)
	// 	b.SendMessage(ctx, &bot.SendMessageParams{
	// 		ChatID:    update.Message.Chat.ID,
	// 		Text:      lexicon.Get("sendError"),
	// 		ParseMode: models.ParseModeHTML,
	// 	})
	// 	return
	// }
	//
	monday := req.Date().AddDate(0, 0, -int(req.Date().Weekday())+1)

	week := resp.IntoWeek()
	res := ""
	res += resp.Group.Name + "\n"
	for _, day := range week.Days {
		newDate := monday.AddDate(0, 0, day.Weekday).Format("02.01.06")
		res += fmt.Sprintf("%d уроков, %s\n", len(day.Lessons), newDate)
		for _, l := range day.Lessons {
			res += strconv.Itoa(l.Index+1) + " " + l.StartTime + " - " + l.EndTime + " " + l.Name + "\n"
		}
		res += "\n\n"
	}

	msg := messages.BuildDayMsg(week.Days[req.Date().Weekday()-1], *req.Date(), resp.Group.Name)

	reply(ctx, b, update, msg)
	logging.Trace("Done handling lesson request")
}
