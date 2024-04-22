package main

import (
	"awesomeProject/internal/bot"
	"awesomeProject/internal/config"
	db2 "awesomeProject/internal/db"
	"awesomeProject/internal/updates_handler"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.LoadEnv()
	doneCh := make(chan os.Signal, 1)

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

	signal.Notify(doneCh, syscall.SIGINT, syscall.SIGTERM)
	<-doneCh
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
}
