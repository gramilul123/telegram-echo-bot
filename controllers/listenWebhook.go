package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gramilul123/telegram-echo-bot/client"
)

func ListenWebhook(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	var update tgbotapi.Update
	json.Unmarshal(bytes, &update)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	bot := client.TgBot{}
	bot.Init()
	bot.Client.Send(msg)
}
