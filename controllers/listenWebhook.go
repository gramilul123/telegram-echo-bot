package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

	exampleQuery := "hello, world"
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:         "Try it 11",
				CallbackData: &exampleQuery,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "Try it 12",
				CallbackData: &exampleQuery,
			},
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:         "Try it 21",
				CallbackData: &exampleQuery,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "Try it 22",
				CallbackData: &exampleQuery,
			},
		),
	)

	if update.Message != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "hello, world!")
		msg.ReplyMarkup = &markup

		bot := client.TgBot{}
		bot.Init()
		bot.Client.Send(msg)
	} else {
		log.Println(update.CallbackQuery.Message.Text)
	}

}
