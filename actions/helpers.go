package actions

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gramilul123/telegram-echo-bot/client"
	"github.com/gramilul123/telegram-echo-bot/db"
	"github.com/gramilul123/telegram-echo-bot/game/strategies"
	war_map "github.com/gramilul123/telegram-echo-bot/game/war_map"
	"github.com/gramilul123/telegram-echo-bot/models"
)

// getMapMsg returns new map
func getNewMapMsg() (text string, gameMap war_map.WarMap) {
	gameMap = war_map.WarMap{}
	gameMap.Create(true)

	text = getTextMap(text, gameMap.Cells)

	return
}

// getTextMap return text view of map
func getTextMap(text string, warCells [][]int) string {

	for i, row := range warCells {
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

	return text
}

// getSelectMapInlineMarkup returns inline buttons
func getSelectMapInlineMarkup() (markup tgbotapi.InlineKeyboardMarkup) {
	buttonList := make(map[string]string)
	buttonList["select_map"] = "Choose map"
	buttonList["accept"] = "Accept"
	markup = getInlineButtons(buttonList)

	return
}

// getEnemyVariants returns inline buttons
func getEnemyVariants() (markup tgbotapi.InlineKeyboardMarkup) {
	buttonList := make(map[string]string)
	buttonList[strategies.SIMPLE] = "Easy"
	buttonList[strategies.MIDDLE] = "Middle"
	markup = getInlineButtons(buttonList)

	return
}

// getNewGameButton return New game inline button
func getNewGameButton() (markup tgbotapi.InlineKeyboardMarkup) {
	buttonList := make(map[string]string)
	buttonList["new"] = "New game"
	markup = getInlineButtons(buttonList)

	return
}

// getNewGameButton return Resign inline button
func getResignButton() (markup tgbotapi.InlineKeyboardMarkup) {

	buttonList := make(map[string]string)
	buttonList["lose"] = "Resign"
	markup = getInlineButtons(buttonList)

	return
}

// getInlineButtons returns inline buttons
func getInlineButtons(buttonList map[string]string) (markup tgbotapi.InlineKeyboardMarkup) {
	brows := [][]tgbotapi.InlineKeyboardButton{}

	brow := []tgbotapi.InlineKeyboardButton{}
	for callback, text := range buttonList {
		callbackSelectMap := callback
		textSelectMap := text
		brow = append(brow, tgbotapi.InlineKeyboardButton{
			Text:         textSelectMap,
			CallbackData: &callbackSelectMap,
		})
	}
	brows = append(brows, brow)
	markup = tgbotapi.NewInlineKeyboardMarkup(brows...)

	return
}

// GetChat returns Chat by chat id
func GetChat(ChatID int64) (chat models.Chat) {
	chats := []models.Chat{}
	models.GetModel(models.CHAT).Get("chat_id", ChatID, &chats)
	chat = chats[0]

	return chat
}

// GetGame returns Game by chat id
func GetGame(ChatID int64) (game models.Game) {
	games := []models.Game{}
	models.GetModel(models.GAME).Get("user_id_one", ChatID, &games)
	game = games[0]

	return game
}

// CreateGame insert row to Game table
func CreateGame(ChatID int64) {
	chat := GetChat(ChatID)

	game := GetGame(ChatID)
	DeleteMessage(ChatID, game.MessageID)
	game = models.Game{
		Status:    models.NewGame,
		UserIDOne: ChatID,
		WarMapOne: chat.AcceptedMap,
	}

	models.GetModel(models.GAME).Delete("user_id_one", ChatID)
	db.Insert(game)
}

// getEmptyWorkMap returns keyboard with emty map
func getEmptyWorkMap() (markup tgbotapi.ReplyKeyboardMarkup) {
	brows := [][]tgbotapi.KeyboardButton{}
	for i := 1; i <= 10; i++ {
		brow := []tgbotapi.KeyboardButton{}
		for j := 1; j <= 10; j++ {
			text := fmt.Sprintf("⬜️\n\n\n%v-%v", i, j)
			brow = append(brow, tgbotapi.KeyboardButton{
				Text: text,
			})
		}
		brows = append(brows, brow)
	}

	markup = tgbotapi.NewReplyKeyboard(brows...)

	return
}

// getWorkMap returns keyboard with work map
func getWorkMap(cells [][]int) (markup tgbotapi.ReplyKeyboardMarkup) {
	var text string
	brows := [][]tgbotapi.KeyboardButton{}
	for i, row := range cells {
		if i > 0 && i < 11 {
			brow := []tgbotapi.KeyboardButton{}
			for j, value := range row {
				if j > 0 && j < 11 {
					switch value {
					case strategies.SHIP:
						text = fmt.Sprint("❌")
					case strategies.HALO:
						text = fmt.Sprint("✴️")
					case strategies.EMPTY:
						text = fmt.Sprint("🔸")
					default:
						text = fmt.Sprintf("⬜️\n\n\n%v-%v", i, j)
					}

					brow = append(brow, tgbotapi.KeyboardButton{
						Text: text,
					})
				}
			}
			brows = append(brows, brow)
		}
	}

	markup = tgbotapi.NewReplyKeyboard(brows...)

	return
}

// getWaitButton returns wait button
func getWaitButton() (markup tgbotapi.ReplyKeyboardMarkup) {
	brows := [][]tgbotapi.KeyboardButton{}
	brow := []tgbotapi.KeyboardButton{}
	brow = append(brow, tgbotapi.KeyboardButton{
		Text: "Wait",
	})
	brows = append(brows, brow)
	markup = tgbotapi.NewReplyKeyboard(brows...)

	return
}

// sendWorkMap send empty work map to keyboard
func sendWorkMap(chatID int64) {
	markup := getEmptyWorkMap()
	msg := tgbotapi.NewMessage(chatID, "Your shot")
	msg.ReplyMarkup = &markup

	response, err := client.Get().Client.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

	game := GetGame(chatID)
	game.MessageID = response.MessageID
	models.GetModel(models.GAME).Update(game, "user_id_one", chatID)
}

// SaveWorkMessageID save work message id to game table
func SaveWorkMessageID(chatID int64, messageID int) {
	game := GetGame(chatID)
	game.MessageID = messageID
	models.GetModel(models.GAME).Update(game, "user_id_one", chatID)
}
