package models

type Chat struct {
	ID      int    `field:"id" key:"primary" extra:"AUTO_INCREMENT"`
	ChatID  int    `field:"chat_id"`
	Command string `field:"command" len:"20"`
	User    string `field:"user" len:"20"`
	Time    int    `field:"time"`
}
