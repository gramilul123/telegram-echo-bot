package routers

import (
	"net/http"

	"github.com/gramilul123/telegram-echo-bot/configs"
	"github.com/gramilul123/telegram-echo-bot/controllers"
)

func Init() {
	http.HandleFunc("/"+configs.GetConfig().Token, func(w http.ResponseWriter, r *http.Request) {
		controllers.SetWebhook(w, r)
	})
	http.HandleFunc("/set_webhook", func(w http.ResponseWriter, r *http.Request) {
		controllers.ListenWebhook(w, r)
	})
}
