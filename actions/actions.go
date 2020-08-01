package actions

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	war_map "github.com/gramilul123/telegram-echo-bot/game/war_map"
	"github.com/gramilul123/telegram-echo-bot/models"
)

// StartBot action after clicking Start button
func StartBot(update tgbotapi.Update) (msg tgbotapi.MessageConfig) {

	chat := models.Chat{
		ChatID: update.Message.Chat.ID,
		Status: models.OPEN,
	}
	models.InstanceChat().DeleteByChatID(chat.ChatID)
	models.InstanceChat().Insert(chat)

	text := getNewMapMsg()
	markup := getSelectMapInlineMarkup()
	msg = tgbotapi.NewMessage(chat.ChatID, text)
	msg.ReplyMarkup = &markup

	return
}

// ReSelectMap action after clicking Select map button
func ReSelectMap(ChatID int64, MessageID int) (editMsg tgbotapi.EditMessageTextConfig) {
	text := getNewMapMsg()
	markup := getSelectMapInlineMarkup()
	editMsg = tgbotapi.NewEditMessageText(ChatID, MessageID, text)
	editMsg.ReplyMarkup = &markup

	return
}

// getMapMsg returns new map
func getNewMapMsg() (text string) {
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

	return
}

// getSelectMapInlineMarkup returns inline buttons
func getSelectMapInlineMarkup() (markup tgbotapi.InlineKeyboardMarkup) {
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

	markup = tgbotapi.NewInlineKeyboardMarkup(brows...)

	return
}
