package main

import (
	"awesomeProject/internal/bot"
	"awesomeProject/internal/config"
	"awesomeProject/internal/ctx"
	db2 "awesomeProject/internal/db"
	"awesomeProject/internal/updates_handler"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.LoadEnv()

	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, syscall.SIGINT, syscall.SIGTERM)

	ctx.InitRootCtx()
	defer func() {
		ctx.GetRootCtx().Cancel()
	}()

	dbConfig := config.GetDBConfig()
	db, err := db2.InitDb(dbConfig)

	if err != nil {
		return
	}
	defer db.Close()

	tgConfig := config.GetTgConfig()
	telegram := bot.NewBot(tgConfig.Token)
	updatesChan := telegram.GetUpdatesChannel()
	go updates_handler.NewConsumer(updatesChan, telegram)

	select {
	case <-doneCh:
		log.Println("finished")
	}
}
