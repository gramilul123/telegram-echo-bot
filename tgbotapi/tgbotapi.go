package tgbotapi

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gramilul123/telegram-echo-bot/configs"
)

type TgBot struct {
	Client *tgbotapi.BotAPI
	Err    error
}

func (bot *TgBot) Init() {
	bot.Client, bot.Err = tgbotapi.NewBotAPI(configs.GetConfig().Token)
	if bot.Err != nil {
		log.Fatal(bot.Err)
	}
	bot.Client.Debug = true
}

func (bot *TgBot) SetWebhook() {
	bot.Init()

	_, bot.Err = bot.Client.SetWebhook(tgbotapi.NewWebhook(fmt.Sprintf("https://%s:%s/%s", configs.GetConfig().URL, configs.GetConfig().Port, configs.GetConfig().Token)))
	if bot.Err != nil {
		log.Fatal(bot.Err)
	}
	bot.GetWebhookInfo()
}

func (bot *TgBot) GetWebhookInfo() {
	_, bot.Err = bot.Client.GetWebhookInfo()
	if bot.Err != nil {
		log.Fatal(bot.Err)
	}
}
