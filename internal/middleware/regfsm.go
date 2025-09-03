package middleware

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"

	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/internal/messages"
	"github.com/madeinheaven91/black-turtle-go/pkg/keyboards"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/models/fsm"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

type FSM struct {
	mu       sync.Mutex
	states   map[int64]fsm.State
	handlers map[fsm.State]fsm.StateHandler
}

func NewFSM() *FSM {
	return &FSM{
		states:   make(map[int64]fsm.State),
		handlers: make(map[fsm.State]fsm.StateHandler),
	}
}

func (f *FSM) SetState(chatID int64, state fsm.State) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.states[chatID] = state
}

func (f *FSM) GetState(chatID int64) (fsm.State, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	state, ok := f.states[chatID]
	return state, ok
}

func (f *FSM) RegisterHandler(state fsm.State, handler func(context.Context, *bot.Bot, *botmodels.Update)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.handlers[state] = handler
}

func (f *FSM) End(chatID int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.states, chatID)
}

func (f *FSM) Middleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		chatID := shared.GetChatID(update)
		if update.CallbackQuery != nil {
			if update.CallbackQuery.Data == "reg_cancel" {
				f.SetState(chatID, fsm.RegCancel)
			}
		}
		currState, ok := f.GetState(chatID)

		if ok {
			if handler, o := f.handlers[currState]; o {
				handler(ctx, b, update)
			} else {
				b.SendMessage(ctx, shared.Params(update, "Unknown state, exited fsm"))
				logging.Error("unknown state %#v")
				f.End(chatID)
			}
			return
		}
		next(ctx, b, update)
	}
}

type FSMHandler struct {
	FSM *FSM
}

func NewFSMHandler(regFSM *FSM) *FSMHandler {
	regFSMHandler := FSMHandler{regFSM}
	regFSM.RegisterHandler(fsm.EnterGroup, regFSMHandler.RegGroupEnter)
	regFSM.RegisterHandler(fsm.EnterTeacher, regFSMHandler.RegTeacherEnter)
	regFSM.RegisterHandler(fsm.RegCancel, regFSMHandler.RegCancel)
	regFSM.RegisterHandler(fsm.MultipleChoice, regFSMHandler.RegMultipleChoice)
	return &regFSMHandler
}

func (fh *FSMHandler) RegGroupEnter(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	chatID := shared.GetChatID(update)
	// Get study entities by name and filter groups
	tmp, err := db.StudyEntities(update.Message.Text)
	studyEntities := make([]models.DBStudyEntity, 0, 1)
	for _, entity := range tmp {
		if entity.Kind == models.Group {
			studyEntities = append(studyEntities, entity)
		}
	}

	if err != nil || len(studyEntities) == 0 {
		logging.Trace("%d unknown group entered (%s)", chatID, update.Message.Text)
		b.SendMessage(ctx, shared.AddReplyMarkup(
			shared.Params(update, lexicon.Get(lexicon.RegGroupNotFound)),
			keyboards.RegCancel(),
		))
	} else if len(studyEntities) == 1 {
		studyEntity := &studyEntities[0]
		err = db.AssignStudyEntity(update, studyEntity)
		if err != nil {
			shared.HandleBotError(err, ctx, b, update)
		} else {
			fh.FSM.End(chatID)
			b.SendMessage(ctx, shared.AddReplyMarkup(
				shared.Params(update, lexicon.Get(lexicon.RegGroupSelected)),
				keyboards.Default(),
			))
			logging.Info("%d registered with %s", chatID, studyEntity.Name)
		}
	} else {
		b.SendMessage(ctx, shared.AddReplyMarkup(
			shared.Params(update, messages.BuildMultipleChoices(models.Group, studyEntities)),
			keyboards.MultipleChoices(models.Group, studyEntities),
		))
		fh.FSM.SetState(chatID, fsm.MultipleChoice)
	}
}

func (fh *FSMHandler) RegTeacherEnter(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	chatID := shared.GetChatID(update)
	tmp, err := db.StudyEntities(update.Message.Text)
	studyEntities := make([]models.DBStudyEntity, 0, 1)
	for _, entity := range tmp {
		if entity.Kind == models.Teacher {
			studyEntities = append(studyEntities, entity)
		}
	}

	if err != nil || len(studyEntities) == 0 {
		logging.Trace("%d unknown teacher entered (%s)", chatID, update.Message.Text)
		b.SendMessage(ctx, shared.AddReplyMarkup(
			shared.Params(update, lexicon.Get(lexicon.RegTeacherNotFound)),
			keyboards.RegCancel(),
		))
	} else if len(studyEntities) == 1 {
		studyEntity := &studyEntities[0]
		err = db.AssignStudyEntity(update, studyEntity)
		if err != nil {
			shared.HandleBotError(err, ctx, b, update)
		} else {
			fh.FSM.End(chatID)
			b.SendMessage(ctx, shared.AddReplyMarkup(
				shared.Params(update, lexicon.Get(lexicon.RegTeacherSelected)),
				keyboards.Default(),
			))
			logging.Info("%d registered with %s", chatID, studyEntity.Name)
		}
	} else {
		b.SendMessage(ctx, shared.AddReplyMarkup(
			shared.Params(update, messages.BuildMultipleChoices(models.Teacher, studyEntities)),
			keyboards.MultipleChoices(models.Teacher, studyEntities),
		))
		fh.FSM.SetState(chatID, fsm.MultipleChoice)
	}
}

func (fh *FSMHandler) RegGroupStart(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	logging.Trace("%d enter group", shared.GetChatID(update))
	fh.FSM.SetState(update.CallbackQuery.Message.Message.Chat.ID, fsm.EnterGroup)
	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      lexicon.Get(lexicon.RegEnterGroup),
		ParseMode: botmodels.ParseModeHTML,
	})
	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: keyboards.RegCancel(),
	})
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
}

func (fh *FSMHandler) RegTeacherStart(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	logging.Trace("%d enter teacher", shared.GetChatID(update))
	fh.FSM.SetState(update.CallbackQuery.Message.Message.Chat.ID, fsm.EnterTeacher)
	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      lexicon.Get(lexicon.RegEnterTeacher),
		ParseMode: botmodels.ParseModeHTML,
	})
	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: keyboards.RegCancel(),
	})
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
}

func (fh *FSMHandler) RegCancel(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	logging.Trace("%d cancelled registration", shared.GetChatID(update))
	fh.FSM.End(update.CallbackQuery.Message.Message.Chat.ID)
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
}

func (fh *FSMHandler) RegMultipleChoice(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	callbackMessage := strings.Split(update.CallbackQuery.Data, "_")
	logging.Trace("%#v", callbackMessage)
	id, err := strconv.Atoi(callbackMessage[2])
	if err != nil {
		logging.Error("incorrect multiple choice message (wrong id): %s", update.CallbackQuery.Data)
		b.SendMessage(ctx, shared.Params(update, lexicon.Error(lexicon.EGeneral)))
		fh.FSM.End(update.CallbackQuery.Message.Message.Chat.ID)
		return
	}
	entity, err := db.StudyEntityByID(id)
	if err != nil {
		logging.Error("incorrect multiple choice message (wrong id): %s", update.CallbackQuery.Data)
		b.SendMessage(ctx, shared.Params(update, lexicon.Error(lexicon.EGeneral)))
		fh.FSM.End(update.CallbackQuery.Message.Message.Chat.ID)
		return
	}
	if err = db.AssignStudyEntity(update, entity); err != nil {
		logging.Error("couldn't assign study entity (%d) to chat %d", entity.ID, update.CallbackQuery.Message.Message.Chat.ID)
		b.SendMessage(ctx, shared.Params(update, lexicon.Error(lexicon.EGeneral)))
		fh.FSM.End(update.CallbackQuery.Message.Message.Chat.ID)
		return
	}
	switch callbackMessage[1] {
	case "group":
		b.SendMessage(ctx, shared.AddReplyMarkup(
			shared.Params(update, lexicon.Get(lexicon.RegGroupSelected)),
			keyboards.Default(),
		))
	case "teacher":
		b.SendMessage(ctx, shared.AddReplyMarkup(
			shared.Params(update, lexicon.Get(lexicon.RegTeacherSelected)),
			keyboards.Default(),
		))
	default:
		logging.Error("incorrect multiple choice message (wrong type): %s", update.CallbackQuery.Data)
		b.SendMessage(ctx, shared.Params(update, lexicon.Error(lexicon.EGeneral)))
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	fh.FSM.End(update.CallbackQuery.Message.Message.Chat.ID)
}
