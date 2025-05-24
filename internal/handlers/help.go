package handlers

import (
	"strings"

	"github.com/go-telegram/bot/models"
)

func HelpMatch(update *models.Update) bool {
	return strings.HasPrefix(strings.ToLower(update.Message.Text), "помощь")
}
