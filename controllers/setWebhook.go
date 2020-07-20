package controllers

import (
	"net/http"

	"github.com/gramilul123/telegram-echo-bot/tgbotapi"
)

func SetWebhook(w http.ResponseWriter, r *http.Request) {
	bot := tgbotapi.TgBot{}

	bot.SetWebhook()
}
