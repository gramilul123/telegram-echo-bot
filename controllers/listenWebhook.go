package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gramilul123/telegram-echo-bot/actions"
	"github.com/gramilul123/telegram-echo-bot/client"
)

// ListenWebhook listens calls from telegram api server
func ListenWebhook(w http.ResponseWriter, r *http.Request) {
	var update tgbotapi.Update
	var msg tgbotapi.MessageConfig

	bytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	json.Unmarshal(bytes, &update)

	log.Println(update.Message.Text)
	if update.Message.Text == "/start" {
		msg = actions.StartBot(update)

	} else if update.Message.Text == "select_map" {
		msg = actions.SelectMap(update)
	}

	client.Get().Client.Send(msg)
}
