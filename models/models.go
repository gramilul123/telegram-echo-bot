package models

import (
	"log"

	"github.com/gramilul123/telegram-echo-bot/db"
)

const (
	CHAT = "Chat"
	GAME = "Game"
)

type Model struct {
	Name     string
	Instance interface{}
}

// GetModel returns instance model
func GetModel(variant string) (model Model) {

	switch variant {
	case GAME:
		model.Instance = &Game{}
	case CHAT:
		model.Instance = &Chat{}
	}

	model.Name = variant

	return
}

// Insert add row to db table
func (Model) Insert(model interface{}) {
	db.Insert(model)
}

// Get returns model by field
func (data Model) Get(field string, value interface{}, models interface{}) interface{} {
	whereValues := make(map[string]interface{})
	whereValues[field] = value
	request := db.GetSelectRequest(data.Name, whereValues)

	err := db.GetDBConnect().DB.Select(models, request)
	if err != nil {
		log.Fatalf("Models: Get: %s", err)
	}

	return models
}

// Update update row by
func (Model) Update(model interface{}, field string, value interface{}) {
	whereValues := make(map[string]interface{})
	whereValues[field] = value
	db.Update(model, whereValues)
}

// Delete removes row by field
func (data Model) Delete(field string, value interface{}) {
	whereValues := make(map[string]interface{})

	whereValues[field] = value
	db.Delete(data.Instance, whereValues)
}
