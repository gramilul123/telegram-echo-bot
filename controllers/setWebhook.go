package controllers

import (
	"github.com/gramilul123/telegram-echo-bot/tgbotapi"
)

func setWebhook() {
	bot := tgbotapi.TgBot{}

	bot.SetWebhook()
}
