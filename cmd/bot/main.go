package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/madeinheaven91/black-turtle-go/internal/handlers"
)

func main() {
	os.Setenv("PUBLICATION_ID", "1c14ffd1-53ea-4c4b-afa5-92b8f80a2183")
	token := os.Getenv("BOT_TOKEN")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {}),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	log.Default().Println("Starting bot...")
	b.RegisterHandlerMatchFunc(handlers.LessonsMatch, handlers.LessonsHandler)
	b.Start(ctx)
}
