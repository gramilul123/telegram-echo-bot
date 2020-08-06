package models

const (
	NewGame    = "N"
	ActiveGame = "A"
	Stop       = "S"
)

type Game struct {
	ID         int    `db:"id" key:"primary" extra:"AUTO_INCREMENT"`
	Status     string `db:"status" len:"1"`
	UserIDOne  int64  `db:"user_id_one"`
	UserIDTwo  int64  `db:"user_id_two"`
	MessageID  int    `db:"message_id"`
	ActiveUser int64  `db:"active_user"`
	WarMapOne  string `db:"war_map_one" type:"text"`
	WarMapTwo  string `db:"war_map_two" type:"text"`
	WorkMapOne string `db:"work_map_one" type:"text"`
	WorkMapTwo string `db:"work_map_two" type:"text"`
}
