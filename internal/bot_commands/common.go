package bot_commands

import (
	"awesomeProject/internal/bot"
	"awesomeProject/internal/db"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/service"
)

var catsRepo *repository.CatRepository = nil
var catsService *service.CatService = nil

func ExtractChatAndMessageToReply(update bot.Update) (chatId int, messageId int) {
	if update.Chat.Id != 0 {
		chatId = update.Chat.Id
	} else {
		chatId = update.From.Id
	}

	return chatId, update.MessageId
}

func GetRepo() *repository.CatRepository {
	if catsRepo == nil {
		catsRepo = repository.NewCatRepository(db.GetDB())
	}

	return catsRepo
}

func GetService() *service.CatService {
	if catsService == nil {
		catsService = service.NewCatService(GetRepo())
	}

	return catsService
}
