package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"

	"github.com/madeinheaven91/black-turtle-go/internal/handlers"
	"github.com/madeinheaven91/black-turtle-go/internal/middleware"
	"github.com/madeinheaven91/black-turtle-go/pkg/config"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/models/fsm"
)

func main() {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		logging.Critical("%s\n", err)
		os.Exit(1)
	}
	for k, v := range envFile {
		os.Setenv(k, v)
	}
	config.InitFromEnv()
	logging.InitLoggers()

	regFSM := middleware.NewFSM()
	regFSMHandler := middleware.NewFSMHandler(regFSM)
	regFSM.RegisterHandler(fsm.EnterGroup, regFSMHandler.RegGroupEnter)
	regFSM.RegisterHandler(fsm.EnterTeacher, regFSMHandler.RegTeacherEnter)
	regFSM.RegisterHandler(fsm.RegCancel, regFSMHandler.RegCancel)

	opts := []bot.Option{
		// empty handler so that stdout isnt being cluttered
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {}),
		bot.WithCallbackQueryDataHandler("start_yes", bot.MatchTypeExact, handlers.StartYes),
		bot.WithCallbackQueryDataHandler("start_no", bot.MatchTypeExact, handlers.StartNo),
		bot.WithCallbackQueryDataHandler("reg_group", bot.MatchTypeExact, regFSMHandler.RegGroupStart),
		bot.WithCallbackQueryDataHandler("reg_teacher", bot.MatchTypeExact, regFSMHandler.RegTeacherStart),
		bot.WithMiddlewares(regFSM.Middleware),
	}

	b, err := bot.New(config.BotToken(), opts...)
	if err != nil {
		logging.Critical("%s\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer func() {
		logging.Info("Recieved interrupt signal, shutting down")
		cancel()
	}()

	logging.Info("Starting bot")
	b.RegisterHandlerMatchFunc(handlers.LessonsMatch, handlers.LessonsHandler, middleware.LogRequest, middleware.DbSync)
	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, handlers.StartHandler, middleware.LogRequest, middleware.DbSync)
	b.RegisterHandlerMatchFunc(handlers.HelpMatch, handlers.HelpHandler, middleware.LogRequest, middleware.DbSync)
	b.RegisterHandlerMatchFunc(handlers.RegistrationMatch, handlers.Registration, middleware.LogRequest, middleware.DbSync)
	b.Start(ctx)
}
