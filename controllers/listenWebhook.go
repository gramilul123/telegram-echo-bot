package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gramilul123/telegram-echo-bot/actions"
	"github.com/gramilul123/telegram-echo-bot/client"
)

// ListenWebhook listens calls from telegram api server
func ListenWebhook(w http.ResponseWriter, r *http.Request) {
	var update tgbotapi.Update
	var msg tgbotapi.MessageConfig
	var editMsg tgbotapi.EditMessageTextConfig

	bytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	json.Unmarshal(bytes, &update)

	if update.Message != nil {
		if update.Message.Text == "/start" {

			msg = actions.StartBot(update)

		} else if update.Message.Text == "select_map" {

			deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
			client.Get().Client.DeleteMessage(deleteMsg)
			msg = actions.SelectMap(update.Message.Chat.ID)

		}
		client.Get().Client.Send(msg)
	} else if update.CallbackQuery != nil {

		if update.CallbackQuery.Data == "select_map" {

			editMsg = actions.ReSelectMap(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
		}

		client.Get().Client.Send(editMsg)

	}

}
