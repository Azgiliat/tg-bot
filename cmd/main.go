package main

import (
	"awesomeProject/internal/bot"
	"awesomeProject/internal/config"
	db2 "awesomeProject/internal/db"
	"awesomeProject/internal/updates_handler"
)

func main() {
	doneCh := make(chan bool)
	config.LoadEnv()

	dbConfig := config.GetDBConfig()
	db, err := db2.InitDb(dbConfig)

	if err != nil {
		return
	}
	defer db.Close()

	tgConfig := config.GetTgConfig()
	telegram := bot.NewBot(tgConfig.Token)
	updatesChan := telegram.GetUpdatesChannel()
	updates_handler.NewConsumer(updatesChan, telegram)

	<-doneCh
}
