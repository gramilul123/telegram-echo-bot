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

	/*exampleQuery := "hello, world"
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:         "1",
				CallbackData: &exampleQuery,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "2",
				CallbackData: &exampleQuery,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "3",
				CallbackData: &exampleQuery,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "4",
				CallbackData: &exampleQuery,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "5",
				CallbackData: &exampleQuery,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "6",
				CallbackData: &exampleQuery,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "7",
				CallbackData: &exampleQuery,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "8",
				CallbackData: &exampleQuery,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "9",
				CallbackData: &exampleQuery,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "10",
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
	)*/
	//markup := tgbotapi.NewInlineKeyboardMarkup()
	/*brows := []tgbotapi.NewInlineKeyboardRow{}
	for i := 1; i <= 10; i++ {
		brow := []tgbotapi.InlineKeyboardButton{}
		for j := 1; j <= 10; j++ {
			callbackText := fmt.Sprintf("%v-%v", i, j)
			brow = append(brow, tgbotapi.InlineKeyboardButton{
				Text:         callbackText,
				CallbackData: &callbackText,
			})
		}
		brows = append(brows, brow)
		//markup = tgbotapi.NewInlineKeyboardMarkup(brow)
	}
	log.Println(brows)
	markup := tgbotapi.NewInlineKeyboardMarkup(brows)*/

	/*btn := tgbotapi.KeyboardButton{
		Text: "",
	}


	markup := tgbotapi.NewReplyKeyboard(
		[]tgbotapi.KeyboardButton{btn, btn, btn, btn, btn, btn, btn, btn, btn, btn},
		[]tgbotapi.KeyboardButton{btn, btn, btn, btn, btn, btn, btn, btn, btn, btn},
		[]tgbotapi.KeyboardButton{btn, btn, btn, btn, btn, btn, btn, btn, btn, btn},
		[]tgbotapi.KeyboardButton{btn, btn, btn, btn, btn, btn, btn, btn, btn, btn},
		[]tgbotapi.KeyboardButton{btn, btn, btn, btn, btn, btn, btn, btn, btn, btn},
		[]tgbotapi.KeyboardButton{btn, btn, btn, btn, btn, btn, btn, btn, btn, btn},
		[]tgbotapi.KeyboardButton{btn, btn, btn, btn, btn, btn, btn, btn, btn, btn},
		[]tgbotapi.KeyboardButton{btn, btn, btn, btn, btn, btn, btn, btn, btn, btn},
		[]tgbotapi.KeyboardButton{btn, btn, btn, btn, btn, btn, btn, btn, btn, btn},
		[]tgbotapi.KeyboardButton{btn, btn, btn, btn, btn, btn, btn, btn, btn, btn})*/

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
