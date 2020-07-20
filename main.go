package main

import (
	"log"
	"net/http"

	"github.com/gramilul123/telegram-echo-bot/configs"
	"github.com/gramilul123/telegram-echo-bot/routers"
)

func main() {
	routers.Init()

	if err := http.ListenAndServe(":"+configs.GetConfig().Port, nil); err != nil {
		log.Fatal(err)
	}
}
