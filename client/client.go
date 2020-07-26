package client

import (
	"fmt"
	"log"
	"sync"

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

	if configs.GetConfig().Debug == "Y" {
		bot.Client.Debug = true
	}
}

func (bot *TgBot) SetWebhook() {
	bot.Init()

	_, bot.Err = bot.Client.SetWebhook(tgbotapi.NewWebhook(fmt.Sprintf("https://%s/%s", configs.GetConfig().URL, configs.GetConfig().Token)))
	if bot.Err != nil {
		log.Fatal(bot.Err)
	}
	bot.GetWebhookInfo()
}

func (bot *TgBot) GetWebhookInfo() tgbotapi.WebhookInfo {
	var info tgbotapi.WebhookInfo

	info, bot.Err = bot.Client.GetWebhookInfo()

	if bot.Err != nil {
		log.Fatal(bot.Err)
	}

	if info.LastErrorDate != 0 {
		log.Fatalf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return info
}

var once sync.Once
var clientTg *TgBot

func Get() *TgBot {
	once.Do(func() {
		clientTg = &TgBot{}
		clientTg.Init()
	})

	return clientTg
}
