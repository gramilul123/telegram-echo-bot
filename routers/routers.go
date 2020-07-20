package routers

import (
	"net/http"

	"github.com/gramilul123/telegram-echo-bot/configs"
)

func Init() {
	http.HandleFunc("/"+configs.GetConfig().Token, controllers.listenWebhook)
	http.HandleFunc("/set_webhook", controllers.setWebhook)
}
