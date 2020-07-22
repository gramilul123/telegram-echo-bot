package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gramilul123/telegram-echo-bot/client"
)

// ListenWebhook listens calls from telegram api server
func ListenWebhook(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	var update tgbotapi.Update
	json.Unmarshal(bytes, &update)

	brows := [][]tgbotapi.KeyboardButton{}
	for i := 1; i <= 10; i++ {
		brow := []tgbotapi.KeyboardButton{}
		for j := 1; j <= 10; j++ {
			text := fmt.Sprintf("%v-%v", i, j)
			brow = append(brow, tgbotapi.KeyboardButton{
				Text: text,
			})
		}
		brows = append(brows, brow)
	}
	markup := tgbotapi.NewReplyKeyboard(brows...)

	if update.Message != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "hello, world!")
		msg.ReplyMarkup = &markup

		bot := client.TgBot{}
		bot.Init()
		bot.Client.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.Text)
		bot := client.TgBot{}
		bot.Init()
		bot.Client.Send(msg)
	}

}
