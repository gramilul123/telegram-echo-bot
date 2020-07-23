package controllers

import (
	"net/http"

	"github.com/gramilul123/telegram-echo-bot/db"
	"github.com/gramilul123/telegram-echo-bot/models"
)

func CreateTable(w http.ResponseWriter, r *http.Request) {
	chat := &models.Chat{}
	db.CreateTable(chat)
}
