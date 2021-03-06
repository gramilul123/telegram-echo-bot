package actions

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
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
func MakeShot(chatID int64, x int, y int) (msg tgbotapi.MessageConfig, result string) {
	var text string
	var markup tgbotapi.ReplyKeyboardMarkup
	WorkMapOne := war_map.WarMap{}
	WarMapTwo := war_map.WarMap{}

	game := GetGame(chatID)

	if game.Status == models.Stop {
		return
	}

	DeleteMessage(chatID, game.MessageID)

	chat := GetChat(chatID)

	if len(game.WorkMapOne) == 0 {
		WorkMapOne.Create(false)
	} else {
		WorkMapOne.JsonToMap(game.WorkMapOne)
	}

	WarMapTwo.JsonToMap(game.WarMapTwo)
	result, WorkMapOne.Cells = strategies.CheckShot(x, y, WorkMapOne.Cells, WarMapTwo)

	if result == strategies.NOK {

		markup = getWaitButton()
		//markup = getWorkMap(WorkMapOne.Cells)
		text = fmt.Sprintf("%d-%d miss", x, y)

	} else if result == strategies.HIT || result == strategies.DESTROYED {

		markup = getWorkMap(WorkMapOne.Cells)

		if result == strategies.HIT {

			text = fmt.Sprintf("%d-%d hit", x, y)

		} else {

			text = fmt.Sprintf("%d-%d Congratulation! Ship destroyed.", x, y)

		}

	} else if result == strategies.WIN {

		editMsg := Finish(chatID, chat.MessageID, "win")
		_, err := client.Get().Client.Send(editMsg)

		if err != nil {
			log.Fatal(err)
		}

		game.Status = models.Stop
		markup = getWorkMap(WorkMapOne.Cells)
		text = fmt.Sprintf("%d-%d Win", x, y)

	} else if result == strategies.DONE {

		log.Fatalln("MakeShot: Unknown error")

	}

	game.WorkMapOne = WorkMapOne.MapToJson()
	game.WarMapTwo = WarMapTwo.MapToJson()
	models.GetModel(models.GAME).Update(game, "user_id_one", chatID)

	if unsafe.Sizeof(markup) != 0 && len(text) > 0 {

		msg = tgbotapi.NewMessage(chat.ChatID, text)
		msg.ReplyMarkup = &markup

	}

	return
}

// DeleteMessage delete message by message id
func DeleteMessage(chatID int64, messageID int) {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)
	client.Get().Client.DeleteMessage(msg)
}

// EnemyGame runs enemy shots
func EnemyGame(chatID int64) {
	var text, result string
	var x, y int
	var cells [][]int

	WarMapOne := war_map.WarMap{}
	game := GetGame(chatID)
	chat := GetChat(chatID)

	str := strategies.GetStrategy(userToEnemy(game.UserIDTwo))

	if len(game.WorkMapTwo) == 0 {
		str.Create()
	} else {
		str.JsonToMap(game.WorkMapTwo)
	}
	WarMapOne.JsonToMap(game.WarMapOne)

	for {
		text += "\n"

		x, y, cells = str.GetShot(result)
		result, cells = strategies.CheckShot(x, y, cells, WarMapOne)

		text = getTextMapWithWork(text, WarMapOne.Cells, cells)

		game.WarMapOne = WarMapOne.MapToJson()
		game.WorkMapTwo = str.MapToJson()
		models.GetModel(models.GAME).Update(game, "user_id_one", chatID)

		if result == strategies.NOK {

			text += fmt.Sprintf("\n\n%d-%d miss", x, y)

		} else if result == strategies.HIT || result == strategies.DESTROYED {

			if result == strategies.HIT {

				text += fmt.Sprintf("\n\n%d-%d hit", x, y)

			} else {

				text += fmt.Sprintf("\n\n%d-%d Congratulation! Ship destroyed.", x, y)

			}

		} else if result == strategies.WIN {

			text += fmt.Sprintf("\n\n%d-%d Win", x, y)

		}

		markup := getResignButton()
		editMsg := tgbotapi.NewEditMessageText(chatID, chat.MessageID, text)
		editMsg.ReplyMarkup = &markup

		_, err := client.Get().Client.Send(editMsg)

		if err != nil {
			log.Fatalf("Listening callback message: %s", err)
		}

		if result == strategies.NOK || result == strategies.WIN {

			if result == strategies.NOK {
				returnUserGame(chatID)
			} else {
				DeleteMessage(chatID, game.MessageID)
			}

			break
		}

		time.Sleep(1 * time.Second)
	}
}
