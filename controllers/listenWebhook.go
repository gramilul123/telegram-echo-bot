package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"unsafe"

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
	var isShot bool
	var result string

	bytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	json.Unmarshal(bytes, &update)

	if update.Message != nil {

		actions.DeleteMessage(update.Message.Chat.ID, update.Message.MessageID)

		if update.Message.Text == "/start" {

			msg, gameMap = actions.StartBot(update)

		} else if update.Message.Text == "Wait" {

		} else if matched, x, y := actions.CheckStep(update.Message.Text); matched {

			msg, result = actions.MakeShot(update.Message.Chat.ID, x, y)
			isShot = true
		}

		if unsafe.Sizeof(msg) != 0 && len(msg.Text) > 0 {
			response, err := client.Get().Client.Send(msg)

			if err != nil {
				log.Fatalf("Listening update message: %s", err)
			}

			if update.Message.Text == "/start" {

				actions.SaveMap(update.Message.Chat.ID, response.MessageID, gameMap)

			} else if isShot {

				actions.SaveWorkMessageID(update.Message.Chat.ID, response.MessageID)

				if result == strategies.NOK {
					actions.EnemyGame(update.Message.Chat.ID)
				}
			}
		}

	} else if update.CallbackQuery != nil {

		if update.CallbackQuery.Data == "select_map" || update.CallbackQuery.Data == "new" {

			editMsg, gameMap = actions.ReSelectMap(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)

			actions.SaveMap(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, gameMap)

			if update.CallbackQuery.Data == "new" {

				game := actions.GetGame(update.CallbackQuery.Message.Chat.ID)
				actions.DeleteMessage(update.CallbackQuery.Message.Chat.ID, game.MessageID)

			}

		} else if update.CallbackQuery.Data == "accept" {

			editMsg = actions.Accept(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)

		} else if update.CallbackQuery.Data == strategies.SIMPLE || update.CallbackQuery.Data == strategies.MIDDLE {

			editMsg = actions.ChoseEnemy(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.Data)

		} else if update.CallbackQuery.Data == "lose" {

			editMsg = actions.Finish(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.Data)

		}

		_, err := client.Get().Client.Send(editMsg)

		if err != nil {
			log.Fatalf("Listening callback message: %s", err)
		}
	}

}
