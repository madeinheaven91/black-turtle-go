package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	mymodels "github.com/madeinheaven91/black-turtle-go/internal/models"
	requests "github.com/madeinheaven91/black-turtle-go/internal/requests"
)

// Send any text message to the bot after the bot has been started

func main() {
	os.Setenv("PUBLICATION_ID", "1c14ffd1-53ea-4c4b-afa5-92b8f80a2183")

	log.Default().Println("Bot started!")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{}



	token := os.Getenv("BOT_TOKEN")
	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "bebra", bot.MatchTypeExact, bebraHandler)
	b.Start(ctx)
}

func bebraHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	date := time.Date(2025, 4, 21, 0, 0, 0, 0, time.UTC)
	text, err := requests.FetchLessons(mymodels.Group, 2, date);
	if err != nil {
		panic(err)
	}

	week := text.IntoWeek()
	res := ""
	for _, day := range week.Days {
		res += strconv.Itoa(day.Weekday) + "\n"
		for _, l := range day.Lessons {
			res += strconv.Itoa(l.Index) + " " + l.StartTime + " - " + l.EndTime + " " + l.Name + "\n"
		}
		res += "\n\n"
	}
	
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: res,
		ParseMode: models.ParseModeHTML,
	})
}
