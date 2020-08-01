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
	ID          int    `db:"id" key:"primary" extra:"AUTO_INCREMENT"`
	ChatID      int64  `db:"chat_id"`
	Status      string `db:"status" len:"1"`
	MessageID   int    `db:"message_id"`
	AcceptedMap string `db:"accepted_map" type:"text"`
}

// DeleteByChatID removes chat by chat id
func (chat Chat) DeleteByChatID(chatID int64) {
	deleteWhere := make(map[string]interface{})

	deleteWhere["chat_id"] = chatID

	db.Delete(chat, deleteWhere)
}

// Insert func inserts row into Chat table
func (Chat) Insert(newChat Chat) {
	db.Insert(newChat)
}

// GetByChatID returns chat object by chat id
func (Chat) GetByChatID(chatID int64) (chat Chat) {
	chats := []Chat{}
	selectWhere := make(map[string]interface{})
	selectWhere["chat_id"] = chatID

	err := db.GetDBConnect().DB.Select(&chats, db.GetSelectRequest(&Chat{}, selectWhere))
	if err != nil {
		log.Fatalf("Chat: GetByChatID: %s", err)
	}

	if len(chats) > 0 {
		chat.ID = chats[0].ID
		chat.ChatID = chats[0].ChatID
		chat.Status = chats[0].Status
		chat.MessageID = chats[0].MessageID
		chat.AcceptedMap = chats[0].AcceptedMap
	}

	return chat
}

// UpdateStatus updates status of chat
func (Chat) UpdateStatus(chat Chat, status string) {
	if chat.ChatID != 0 {
		chat.Status = status
		db.UpdateRow(chat, "chat_id")
	}
}

// InstanceChat returns chat instans
func InstanceChat() Chat {
	chat := Chat{}

	return chat
}
