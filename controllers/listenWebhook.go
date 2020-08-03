package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gramilul123/telegram-echo-bot/actions"
	"github.com/gramilul123/telegram-echo-bot/client"
	"github.com/gramilul123/telegram-echo-bot/game/strategies"
	"github.com/gramilul123/telegram-echo-bot/game/war_map"
)

// ListenWebhook listens calls from telegram api server
func ListenWebhook(w http.ResponseWriter, r *http.Request) {
	var update tgbotapi.Update
	var msg tgbotapi.MessageConfig
	var editMsg tgbotapi.EditMessageTextConfig
	var gameMap war_map.WarMap

	bytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	json.Unmarshal(bytes, &update)

	if update.Message != nil {

		if update.Message.Text == "/start" {

			msg, gameMap = actions.StartBot(update)

			response, err := client.Get().Client.Send(msg)

			if err != nil {
				log.Fatal(err)
			}

			actions.SaveMap(update.Message.Chat.ID, response.MessageID, gameMap)
		}
	} else if update.CallbackQuery != nil {

		if update.CallbackQuery.Data == "select_map" {

			editMsg, gameMap = actions.ReSelectMap(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)

			_, err := client.Get().Client.Send(editMsg)

			if err != nil {
				log.Fatal(err)
			}

			actions.SaveMap(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, gameMap)

		} else if update.CallbackQuery.Data == "accept" {

			editMsg = actions.Accept(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)

			_, err := client.Get().Client.Send(editMsg)

			if err != nil {
				log.Fatal(err)
			}

		} else if update.CallbackQuery.Data == strategies.SIMPLE || update.CallbackQuery.Data == strategies.MIDDLE {

			editMsg = actions.ChoseEnemy(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.Data)

			_, err := client.Get().Client.Send(editMsg)

			if err != nil {
				log.Fatal(err)
			}

		} else if update.CallbackQuery.Data == "lose" {

			editMsg = actions.Finish(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.Data)

			_, err := client.Get().Client.Send(editMsg)

			if err != nil {
				log.Fatal(err)
			}
		}

	}

}
