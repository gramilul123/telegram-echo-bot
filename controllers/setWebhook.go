package controllers

import (
	"log"
	"net/http"

	"github.com/gramilul123/telegram-echo-bot/client"
)

func SetWebhook(w http.ResponseWriter, r *http.Request) {
	bot := client.TgBot{}
	bot.SetWebhook()
	info := bot.GetWebhookInfo()

	log.Println(info)
}
