package shared

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/errors"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
)

func GetChatName(update *botmodels.Update) string {
	var name string
	if update.Message.Chat.Type == "private" {
		if update.Message.Chat.FirstName != "" && update.Message.Chat.LastName != "" {
			name = update.Message.Chat.FirstName + " " + update.Message.Chat.LastName
		} else if update.Message.Chat.FirstName != "" {
			name = update.Message.Chat.FirstName
		} else {
			name = "<unknown>"
		}
	} else {
		if update.Message.Chat.Title != "" {
			name = update.Message.Chat.Title
		} else {
			name = "<unknown>"
		}
	}
	return name
}

func GetFromName(update *botmodels.Update) string {
	var name string
	if update.Message.From.FirstName != "" && update.Message.From.LastName != "" {
		name = update.Message.From.FirstName + " " + update.Message.From.LastName
	} else if update.Message.From.FirstName != "" {
		name = update.Message.From.FirstName
	} else if update.Message.Chat.Title != "" {
		name = update.Message.Chat.Title
	} else {
		name = "<unknown>"
	}
	return name
}

// Normalizes Go Weekday type. Sunday = 0, ..., Saturday = 6 becomes Monday = 0, ..., Sunday = 6
func NormalizeWeekday(weekday time.Weekday) int {
	return (int(weekday) + 6) % 7
}

// Returns monday of the week that the input date belongs to
func GetMonday(date time.Time) time.Time {
	return date.AddDate(0, 0, -int(NormalizeWeekday(date.Weekday())))
}

func HandleBotError(err error, ctx context.Context, b *bot.Bot, update *botmodels.Update) bool {
	if err != nil {
		e, ok := err.(*errors.BotError)
		if ok {
			logging.Error("%q", e)
			b.SendMessage(ctx, Params(update, e.Display()))
		} else {
			logging.Error("%q", err)
			b.SendMessage(ctx, Params(update, lexicon.Error(lexicon.EGeneral)))
		}
		return false
	}
	return true
}

func Params(update *botmodels.Update, text string) *bot.SendMessageParams {
	return &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: text,
		ParseMode: botmodels.ParseModeHTML,
	}

}

func AddReplyMarkup(params *bot.SendMessageParams, kb botmodels.ReplyMarkup) *bot.SendMessageParams {
	params.ReplyMarkup = kb
	return params
}

func GetChatID(update *botmodels.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}
	if update.CallbackQuery.Message.Message != nil {
		return update.CallbackQuery.Message.Message.Chat.ID
	}
	return 0
}
