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


	regCancel = botmodels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botmodels.InlineKeyboardButton{
			{
				{Text: "Отмена", CallbackData: "reg_cancel" },
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

	start = botmodels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botmodels.InlineKeyboardButton{
			{
				{Text: "Да, давай!", CallbackData: "start_yes"},
				{Text: "Нет, спасибо", CallbackData: "start_no"},
			},
		},
	}
	chooseGroupOrTeacher = botmodels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botmodels.InlineKeyboardButton{
			{
				{Text: "Группы", CallbackData: "reg_group"},
				{Text: "Преподавателя", CallbackData: "reg_teacher"},
			},
		},
	}
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

func Start() *botmodels.InlineKeyboardMarkup {
	return &start
}

func ChooseGroupOrTeacher() *botmodels.InlineKeyboardMarkup {
	return &chooseGroupOrTeacher
}

func RegCancel() *botmodels.InlineKeyboardMarkup {
	return &regCancel
}
