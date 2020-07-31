package models

import (
	"log"

	"github.com/gramilul123/telegram-echo-bot/db"
)

const (
	OPEN = "O" //Chat open
	WAIT = "W" //Chat waiting
)

type Chat struct {
	ID      int    `db:"id" key:"primary" extra:"AUTO_INCREMENT"`
	ChatID  int64  `db:"chat_id"`
	Command string `db:"command" len:"20"`
	User    string `db:"user" len:"20"`
	Status  string `db:"status" len:"1"`
	Time    int    `db:"time"`
}

func (chat Chat) DeleteByChatID(chatId int64) {

	deleteWhere := []string{}
	deleteWhere = append(deleteWhere, "chat_id")

	db.Delete(chat, deleteWhere)
}

// Insert func inserts row into Chat table
func (Chat) Insert(chat Chat) {
	db.Insert(chat)
}

func (Chat) GetByChatId(chatId int64) Chat {
	chats := []Chat{}
	selectWhere := make(map[string]interface{})
	selectWhere["chat_id"] = chatId

	err := db.GetDBConnect().DB.Select(&chats, db.GetSelectRequest(&Chat{}, selectWhere))
	if err != nil {
		log.Fatalf("Chat: GetByChatId: %s", err)
	}

	if len(chats) == 0 {
		log.Fatal("Chat: GetByChatId: Chat %d not found", chatId)
	}

	return chats[0]
}

func (Chat) UpdateStatus(chat Chat, status string) Chat {
	chat.Status = status
	db.UpdateRow(chat, "chat_id")

	return chat
}

func Chats(chatID ...int64) *Chat {
	chat := &Chat{}

	return chat
}
