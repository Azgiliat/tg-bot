package bot_commands

import (
	"awesomeProject/internal/bot"
	"errors"
)

func GetCat(update bot.Update, cat string) (int, int, string, error) {
	chatId, messageId := ExtractChatAndMessageToReply(update)

	if cat == "" {
		return chatId, messageId, "", errors.New("cat is empty string")
	}

	catPhoto := GetService().GetCatPhotosByTag(cat)

	if catPhoto != nil {
		return chatId, messageId, catPhoto.Link, nil
	} else {
		return chatId, messageId, "", errors.New("error during get cats by tag")
	}
}
