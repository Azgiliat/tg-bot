package bot_commands

import (
	"awesomeProject/internal/bot"
)

// GetCats TODO add cache for available cats, check if cat exist not in db but here
func GetCats(update bot.Update) (int, []string) {
	var idToReply int

	if update.Chat.Id != 0 {
		idToReply = update.Chat.Id
	} else {
		idToReply = update.From.Id
	}

	cats := GetService().GetAvailableCats()

	return idToReply, cats
}
