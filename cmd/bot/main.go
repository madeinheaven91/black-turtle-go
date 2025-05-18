package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"

	"github.com/madeinheaven91/black-turtle-go/internal/handlers"
	"github.com/madeinheaven91/black-turtle-go/internal/logging"
	"github.com/madeinheaven91/black-turtle-go/internal/middleware"
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
	logging.InitLogLevel()


	opts := []bot.Option{
		// empty handler so that stdout isnt being cluttered
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {}),
	}

	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
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
	b.Start(ctx)
}
