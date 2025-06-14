package middleware

import (
	"context"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"

	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func LogRequest(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		name := shared.GetFromName(update)
		if update.Message.Chat.Type == "private" {
			logging.Telegram("%s (%d): %s\n", name, update.Message.Chat.ID, update.Message.Text)
		} else {
			chat_name := shared.GetChatName(update)
			logging.Telegram("%s (%d) in %s (%d): %s\n", name, update.Message.From.ID, chat_name, update.Message.Chat.ID, update.Message.Text)
		}
		next(ctx, b, update)
	}
}

type pair struct {
	fst bool
	snd bool
}

func DbSync(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		conn := db.GetConnection()
		defer db.CloseConn(conn)
		row := conn.QueryRow(context.Background(), "select * from chat where id=$1", update.Message.Chat.ID)
		var name string
		var username string
		err := row.Scan(nil, nil, &name, &username, nil, nil)
		if err != nil {
			logging.Error("%s\n", err)
			next(ctx, b, update)
			return
		}

		newName := shared.GetChatName(update)
		newUsername := update.Message.Chat.Username
		matches := pair{
			name == newName,
			username == newUsername,
		}
		switch matches {
		case pair{false, true}:
			_, err = conn.Exec(context.Background(), "update chat set name=$1 where id=$2", newName, update.Message.Chat.ID)
		case pair{true, false}:
			_, err = conn.Exec(context.Background(), "update chat set username=$1 where id=$2", newUsername, update.Message.Chat.ID)
		case pair{false, false}:
			_, err = conn.Exec(context.Background(), "update chat set name=$1 username=$1 where id=$2", newName, newUsername, update.Message.Chat.ID)
		}
		if err != nil {
			logging.Error("%s\n", err)
		}
		next(ctx, b, update)
	}
}
