package actions

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gramilul123/telegram-echo-bot/models"
)

func StartBot(update tgbotapi.Update) (msg tgbotapi.MessageConfig) {

	chat := models.Chat{
		ChatID:  update.Message.Chat.ID,
		Command: update.Message.Text,
		Status:  models.OPEN,
	}
	models.Chats().DeleteByChatID(chat.ChatID)
	models.Chats().Insert(chat)
	msg = tgbotapi.NewMessage(chat.ChatID, "Hello!")

	brows := [][]tgbotapi.KeyboardButton{}
	brow := []tgbotapi.KeyboardButton{}
	text := fmt.Sprintf("select_map")
	brow = append(brow, tgbotapi.KeyboardButton{
		Text: text,
	})
	brows = append(brows, brow)
	markup := tgbotapi.NewReplyKeyboard(brows...)
	msg.ReplyMarkup = &markup

	return
}

func SelectMap(update tgbotapi.Update) (msg tgbotapi.MessageConfig) {
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Your map")
	brows := [][]tgbotapi.KeyboardButton{}
	for i := 1; i <= 10; i++ {
		brow := []tgbotapi.KeyboardButton{}
		for j := 1; j <= 10; j++ {
			text := fmt.Sprintf("⬜️ \n \n \n \n %d-%d", i, j)
			brow = append(brow, tgbotapi.KeyboardButton{
				Text: text,
			})
		}
		brows = append(brows, brow)
	}
	markup := tgbotapi.NewReplyKeyboard(brows...)
	msg.ReplyMarkup = &markup

	return
}
