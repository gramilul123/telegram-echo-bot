package actions

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"unsafe"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gramilul123/telegram-echo-bot/client"
	"github.com/gramilul123/telegram-echo-bot/db"
	"github.com/gramilul123/telegram-echo-bot/game/strategies"
	war_map "github.com/gramilul123/telegram-echo-bot/game/war_map"
	"github.com/gramilul123/telegram-echo-bot/models"
)

var EnemyToUser = map[string]int64{
	strategies.SIMPLE: 1,
	strategies.MIDDLE: 2,
}

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
	chat := GetChat(ChatID)

	chat.MessageID = MessageID
	chat.AcceptedMap = gameMap.MapToJson()

	models.GetModel(models.CHAT).Update(chat, "chat_id", ChatID)
}

// Accept action after clicking Accept map
func Accept(ChatID int64, MessageID int) (editMsg tgbotapi.EditMessageTextConfig) {
	text := "Choose enemy"
	markup := getEnemyVariants()
	editMsg = tgbotapi.NewEditMessageText(ChatID, MessageID, text)
	editMsg.ReplyMarkup = &markup

	CreateGame(ChatID)

	return
}

// Finish action
func Finish(ChatID int64, MessageID int, result string) (editMsg tgbotapi.EditMessageTextConfig) {
	text := result
	markup := getNewGameButton()
	editMsg = tgbotapi.NewEditMessageText(ChatID, MessageID, text)
	editMsg.ReplyMarkup = &markup

	return
}

// ChoseEnemy action after clicking enemy variant
func ChoseEnemy(ChatID int64, MessageID int, enemy string) (editMsg tgbotapi.EditMessageTextConfig) {
	var text string

	game := GetGame(ChatID)

	gameMapTwo := war_map.WarMap{}
	gameMapTwo.Create(true)
	game.UserIDTwo = EnemyToUser[enemy]
	game.WarMapTwo = gameMapTwo.MapToJson()

	models.GetModel(models.GAME).Update(game, "user_id_one", ChatID)

	gameMapOne := war_map.WarMap{}
	gameMapOne.JsonToMap(game.WarMapOne)

	text = getTextMap(text, gameMapOne.Cells)

	markup := getResignButton()
	editMsg = tgbotapi.NewEditMessageText(ChatID, MessageID, text)
	editMsg.ReplyMarkup = &markup

	sendWorkMap(ChatID)

	return
}

// CheckStep returns result checking
func CheckStep(text string) (matched bool, x int, y int) {
	rexp, _ := regexp.Compile(`(\d+)-(\d+)`)
	str := rexp.FindStringSubmatch(text)

	if len(str) == 3 {
		matched = true
		x, _ = strconv.Atoi(str[1])
		y, _ = strconv.Atoi(str[2])
	}

	return
}

// MakeShot return result user shot
func MakeShot(chatID int64, x int, y int) (msg tgbotapi.MessageConfig) {
	var result, text string
	var markup tgbotapi.ReplyKeyboardMarkup
	var editMsg tgbotapi.EditMessageTextConfig
	WorkMapOne := war_map.WarMap{}
	WarMapTwo := war_map.WarMap{}

	game := GetGame(chatID)

	chat := GetChat(chatID)

	if len(game.WorkMapOne) == 0 {
		WorkMapOne.Create(false)
	} else {
		WorkMapOne.JsonToMap(game.WorkMapOne)
	}

	WarMapTwo.JsonToMap(game.WarMapTwo)
	result, WorkMapOne.Cells = strategies.CheckShot(x, y, WorkMapOne.Cells, WarMapTwo)
	log.Println(result)
	if result == strategies.NOK {

		//markup = getWaitButton()
		markup = getWorkMap(WorkMapOne.Cells)
		text = fmt.Sprintf("%d-%d miss", x, y)

	} else if result == strategies.HIT || result == strategies.DESTROYED {

		markup = getWorkMap(WorkMapOne.Cells)
		text = fmt.Sprintf("%d-%d hit", x, y)

	} else if result == strategies.WIN {

		editMsg = Finish(chatID, chat.MessageID, "win")
		_, err := client.Get().Client.Send(editMsg)

		if err != nil {
			log.Fatal(err)
		}

	} else if result == strategies.DONE {

		log.Fatalln("MakeShot: Unknown error")

	}

	game.WorkMapOne = WorkMapOne.MapToJson()
	game.WarMapTwo = WarMapTwo.MapToJson()
	models.GetModel(models.GAME).Update(game, "user_id_one", chatID)

	if unsafe.Sizeof(markup) != 0 {

		msg = tgbotapi.NewMessage(chat.ChatID, text)
		msg.ReplyMarkup = &markup

	} else if unsafe.Sizeof(editMsg) != 0 {

	}

	return
}
