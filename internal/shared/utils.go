package shared

import (
	botmodels "github.com/go-telegram/bot/models"
)

func GetChatName(update *botmodels.Update) string {
	var name string
	if update.Message.Chat.Type == "private" {
		if update.Message.Chat.FirstName != "" && update.Message.Chat.LastName != "" {
			name = update.Message.Chat.FirstName + " " + update.Message.Chat.LastName
		} else if update.Message.Chat.FirstName != "" {
			name = update.Message.Chat.FirstName
		} else {
			name = "<unknown>"
		}
	} else {
		if update.Message.Chat.Title != "" {
			name = update.Message.Chat.Title
		} else {
			name = "<unknown>"
		}
	}
	return name
}

func GetFromName(update *botmodels.Update) string {
	var name string
	if update.Message.From.FirstName != "" && update.Message.From.LastName != "" {
		name = update.Message.From.FirstName + " " + update.Message.From.LastName
	} else if update.Message.From.FirstName != "" {
		name = update.Message.From.FirstName
	} else if update.Message.Chat.Title != "" {
		name = update.Message.Chat.Title
	} else {
		name = "<unknown>"
	}
	return name
}
