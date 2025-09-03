package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"

	"github.com/madeinheaven91/black-turtle-go/internal/handlers"
	"github.com/madeinheaven91/black-turtle-go/internal/handlers/admin"
	mw "github.com/madeinheaven91/black-turtle-go/internal/middleware"
	"github.com/madeinheaven91/black-turtle-go/pkg/config"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
)

func main() {
	// Reading env variables and initializing config and logging
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

	// Initializing registration FSM
	regFSM := mw.NewFSM()
	regFSMHandler := mw.NewFSMHandler(regFSM)

	opts := []bot.Option{
		// empty handler so that stdout isnt being cluttered
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {}),
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

	container := handlers.NewContainer(b, admin.NewNotificationService(b, 10))

	logging.Info("Starting bot")

	mwChain := func(next bot.HandlerFunc) bot.HandlerFunc {
		return mw.LogRequest(mw.DBSync(mw.Recover(next)))
	}

	// Registering handlers
	b.RegisterHandlerMatchFunc(handlers.LessonsMatch, container.LessonsHandler, mwChain)
	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, container.StartHandler, mwChain)
	b.RegisterHandlerMatchFunc(handlers.HelpMatch, container.HelpHandler, mwChain)
	b.RegisterHandlerMatchFunc(handlers.FioMatch, container.FioHandler, mw.DBSync, mw.LogRequest)

	b.RegisterHandlerMatchFunc(handlers.RegistrationMatch, container.Registration, mw.DBSync, mw.LogRequest)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "start_yes", bot.MatchTypeExact, handlers.StartYes)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "start_no", bot.MatchTypeExact, handlers.StartNo)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "reg_group", bot.MatchTypeExact, regFSMHandler.RegGroupStart)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "reg_teacher", bot.MatchTypeExact, regFSMHandler.RegTeacherStart)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "choose", bot.MatchTypePrefix, regFSMHandler.RegMultipleChoice)

	b.RegisterHandler(bot.HandlerTypeMessageText, "send", bot.MatchTypeCommand, container.SendHandler, mwChain, mw.CheckAdmin)
	b.RegisterHandler(bot.HandlerTypeMessageText, "stat", bot.MatchTypeCommand, container.StatHandler, mwChain, mw.CheckAdmin)

	b.Start(ctx)
}
