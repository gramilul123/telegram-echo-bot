package models

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
