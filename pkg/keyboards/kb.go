package keyboards

import (
	"fmt"
	"strconv"

	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/models"
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
				{Text: "Отмена", CallbackData: "reg_cancel"},
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

func MultipleChoices(kind models.StudyEntityType, choices []models.DBStudyEntity) *botmodels.InlineKeyboardMarkup {
	// NOTE: intended to be used only with multiple choice message builder, so
	// the order of entities is meant to be preserved. Thats why there are no additional
	// checks
	keyboard := make([][]botmodels.InlineKeyboardButton, 0, 1)
	row := 0
	for i, choice := range choices {
		if len(keyboard) <= row {
			keyboard = append(keyboard, make([]botmodels.InlineKeyboardButton, 0))
		}
		keyboard[row] = append(keyboard[row], botmodels.InlineKeyboardButton{
			Text:         strconv.Itoa(i + 1),
			CallbackData: fmt.Sprintf("choose_%s_%d", kind, choice.ID),
		})
		if len(keyboard[row]) == 8 {
			row++
		}
	}
	res := botmodels.InlineKeyboardMarkup{InlineKeyboard: keyboard}
	return &res
}
