package keyboards

import (
	botmodels "github.com/go-telegram/bot/models"
)

var (
	empty = botmodels.ReplyKeyboardMarkup{}
	def   = botmodels.ReplyKeyboardMarkup{
		IsPersistent:   true,
		ResizeKeyboard: true,
		Keyboard: [][]botmodels.KeyboardButton{
			{
				{Text: "Пары"},
				{Text: "Пары завтра"},
			},
			{
				{Text: "Помощь"},
			},
		},
	}

	help = botmodels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botmodels.InlineKeyboardButton{
			{
				{Text: "Техподдержка", URL: "tg://user?id=2087648271"},
			},
		},
	}
	chooseGroupOrTeacher = botmodels.InlineKeyboardMarkup{}
)

func Empty() *botmodels.ReplyKeyboardMarkup {
	return &empty
}

func Default() *botmodels.ReplyKeyboardMarkup {
	return &def
}

func Help() *botmodels.InlineKeyboardMarkup {
	return &help
}
