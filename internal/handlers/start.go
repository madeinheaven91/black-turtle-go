package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/pkg/keyboards"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func (c Container) StartHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	chat := db.Chat(update.Message.Chat.ID)
	if chat != nil {
		logging.Info("Chat (%d) already exists", chat.ID)
	} else {
		logging.Info("New chat (%d)!", update.Message.Chat.ID)
		db.AddChat(update)
	}
	b.SendMessage(ctx, shared.AddReplyMarkup(shared.Params(update, lexicon.Get(lexicon.Start)), keyboards.Start()))
}

func StartNo(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      lexicon.Get(lexicon.RegCancel),
		ParseMode: botmodels.ParseModeHTML,
	})
	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: keyboards.Empty(),
	})
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	logging.Trace("Chat (%d) cancelled registration", update.CallbackQuery.Message.Message.Chat.ID)
}

func StartYes(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      lexicon.Get(lexicon.RegGroupOrTeacher),
		ParseMode: botmodels.ParseModeHTML,
	})
	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: keyboards.ChooseGroupOrTeacher(),
	})
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
}
