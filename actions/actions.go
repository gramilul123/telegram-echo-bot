package actions

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	war_map "github.com/gramilul123/telegram-echo-bot/game/war_map"
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
	brow = append(brow, tgbotapi.KeyboardButton{
		Text: "select_map",
	})
	brows = append(brows, brow)
	markup := tgbotapi.NewReplyKeyboard(brows...)
	msg.ReplyMarkup = &markup

	return
}

func SelectMap(ChatID int64) (msg tgbotapi.MessageConfig) {
	text := "Your map:\n"
	gameMap := war_map.WarMap{}
	gameMap.Create(true)

	for i, row := range gameMap.Cells {
		if i > 0 && i < 11 {
			for j, cell := range row {
				if j > 0 && j < 11 {
					if cell == war_map.Ship {
						text += "⬛️"
					} else {
						text += "⬜️"
					}
				}
			}
			text += "\n"
		}
	}

	msg = tgbotapi.NewMessage(ChatID, text)

	brows := [][]tgbotapi.InlineKeyboardButton{}

	brow := []tgbotapi.InlineKeyboardButton{}

	callbackSelectMap := fmt.Sprintf("select_map")
	textSelectMap := fmt.Sprintf("Select map")
	brow = append(brow, tgbotapi.InlineKeyboardButton{
		Text:         textSelectMap,
		CallbackData: &callbackSelectMap,
	})

	callbackAccept := fmt.Sprintf("accept")
	textAccept := fmt.Sprintf("Accept")
	brow = append(brow, tgbotapi.InlineKeyboardButton{
		Text:         textAccept,
		CallbackData: &callbackAccept,
	})

	brows = append(brows, brow)

	markup := tgbotapi.NewInlineKeyboardMarkup(brows...)

	msg.ReplyMarkup = &markup

	return
}

func ReSelectMap(ChatID int64, MessageID int) (editMsg tgbotapi.EditMessageTextConfig) {
	text := "Your map:\n"
	gameMap := war_map.WarMap{}
	gameMap.Create(true)

	for i, row := range gameMap.Cells {
		if i > 0 && i < 11 {
			for j, cell := range row {
				if j > 0 && j < 11 {
					if cell == war_map.Ship {
						text += "⬛️"
					} else {
						text += "⬜️"
					}
				}
			}
			text += "\n"
		}
	}
	editMsg = tgbotapi.NewEditMessageText(ChatID, MessageID, text)

	brows := [][]tgbotapi.InlineKeyboardButton{}

	brow := []tgbotapi.InlineKeyboardButton{}

	callbackSelectMap := fmt.Sprintf("select_map")
	textSelectMap := fmt.Sprintf("Select map")
	brow = append(brow, tgbotapi.InlineKeyboardButton{
		Text:         textSelectMap,
		CallbackData: &callbackSelectMap,
	})

	callbackAccept := fmt.Sprintf("accept")
	textAccept := fmt.Sprintf("Accept")
	brow = append(brow, tgbotapi.InlineKeyboardButton{
		Text:         textAccept,
		CallbackData: &callbackAccept,
	})

	brows = append(brows, brow)

	markup := tgbotapi.NewInlineKeyboardMarkup(brows...)

	editMsg.ReplyMarkup = &markup

	return
}
