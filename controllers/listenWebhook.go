package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	botclien "github.com/gramilul123/telegram-echo-bot/tgbotapi"
)

func listenWebhook(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	var update tgbotapi.Update
	json.Unmarshal(bytes, &update)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	bot := botclien.TgBot{}.Client
	bot.Send(msg)
}
