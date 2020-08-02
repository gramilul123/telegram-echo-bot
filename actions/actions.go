package actions

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gramilul123/telegram-echo-bot/db"
	"github.com/gramilul123/telegram-echo-bot/game/strategies"
	war_map "github.com/gramilul123/telegram-echo-bot/game/war_map"
	"github.com/gramilul123/telegram-echo-bot/models"
)

// StartBot action after clicking Start button
func StartBot(update tgbotapi.Update) (msg tgbotapi.MessageConfig, gameMap war_map.WarMap) {
	chat := models.Chat{
		ChatID: update.Message.Chat.ID,
		Status: models.OPEN,
	}

	models.GetModel(models.CHAT).Delete("chat_id", chat.ID)
	db.Insert(chat)

	text, gameMap := getNewMapMsg()
	markup := getSelectMapInlineMarkup()
	msg = tgbotapi.NewMessage(chat.ChatID, text)
	msg.ReplyMarkup = &markup

	return
}

// ReSelectMap action after clicking Select map button
func ReSelectMap(ChatID int64, MessageID int) (editMsg tgbotapi.EditMessageTextConfig, gameMap war_map.WarMap) {
	text, gameMap := getNewMapMsg()
	markup := getSelectMapInlineMarkup()
	editMsg = tgbotapi.NewEditMessageText(ChatID, MessageID, text)
	editMsg.ReplyMarkup = &markup

	return
}

// SaveMap save data to chat table
func SaveMap(ChatID int64, MessageID int, gameMap war_map.WarMap) {
	chats := []models.Chat{}

	models.GetModel(models.CHAT).Get("chat_id", ChatID, &chats)
	chat := chats[0]

	chat.MessageID = MessageID
	chat.AcceptedMap = gameMap.MapToJson()

	models.GetModel(models.CHAT).Update(chat, "chat_id", ChatID)
}

func CreateGame(ChatID int64) {
	chats := []models.Chat{}

	models.GetModel(models.CHAT).Get("chat_id", ChatID, &chats)
	chat := chats[0]

	game := models.Game{
		Status:    models.NewGame,
		UserIDOne: ChatID,
		WarMapOne: chat.AcceptedMap,
	}

	models.GetModel(models.GAME).Delete("user_id_one", ChatID)
	db.Insert(game)
}

// Accept action after clicking Accept map
func Accept(ChatID int64, MessageID int) (editMsg tgbotapi.EditMessageTextConfig) {
	text := "Choose enemy"
	brows := [][]tgbotapi.InlineKeyboardButton{}
	brow := []tgbotapi.InlineKeyboardButton{}

	callbackSelectMap := fmt.Sprintf(strategies.SIMPLE)
	textSelectMap := fmt.Sprintf("Easy")
	brow = append(brow, tgbotapi.InlineKeyboardButton{
		Text:         textSelectMap,
		CallbackData: &callbackSelectMap,
	})

	callbackAccept := fmt.Sprintf(strategies.MIDDLE)
	textAccept := fmt.Sprintf("Middle")
	brow = append(brow, tgbotapi.InlineKeyboardButton{
		Text:         textAccept,
		CallbackData: &callbackAccept,
	})

	brows = append(brows, brow)

	markup := tgbotapi.NewInlineKeyboardMarkup(brows...)
	editMsg = tgbotapi.NewEditMessageText(ChatID, MessageID, text)
	editMsg.ReplyMarkup = &markup

	CreateGame(ChatID)

	return
}

// getMapMsg returns new map
func getNewMapMsg() (text string, gameMap war_map.WarMap) {
	gameMap = war_map.WarMap{}
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
