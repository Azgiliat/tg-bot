package updates_handler

import (
  "awesomeProject/internal/bot"
  botcommands2 "awesomeProject/internal/bot_commands"
  "fmt"
  "log"
  "regexp"
  "strings"
)

func NewHandler(ch <-chan bot.Update, telegram *bot.TgBot) {
  rg, err := regexp.Compile("/([a-z_]+) ?([a-zA-Z]+)?")

  if err != nil {
    return
  }

  for update := range ch {
    var isPhotoAttachedToMessage bool = len(update.Message.Photo) != 0
    var entities *[]bot.Entity
    var textToParse *string
    var messageWithPhoto *bot.Message

    if isPhotoAttachedToMessage {
      entities = &update.CaptionEntities
      textToParse = &update.Caption
      messageWithPhoto = &update.Message
    } else if len(update.Entities) != 0 {
      entities = &update.Entities
      textToParse = &update.Text
      messageWithPhoto = update.ReplyToMessage
    }

    if entities == nil || textToParse == nil || len(*entities) == 0 {
      log.Println("empty entities or command")
      continue
    }

    for _, entity := range *entities {
      if entity.Type == "bot_command" {
        parsedMessage := rg.FindStringSubmatch(*textToParse)

        if len(parsedMessage) < 3 {
          continue
        }

        command := parsedMessage[1]

        switch command {
        case "store_cat":
          go func() {
            var err error
            var responseText string

            photoId := botcommands2.ExtractPhotoId(messageWithPhoto)

            if photoId == "" {
              log.Println("failed to extract photo id")
              return
            }

            res, err := telegram.DownloadFile(photoId)

            if err != nil {
              return
            }

            err = botcommands2.StoreCat(parsedMessage[2], res)
            chatId, messageId := botcommands2.ExtractChatAndMessageToReply(update)

            if err != nil {
              log.Println(err)
              responseText = "Failed to store cat"
            } else {
              responseText = "Cat has been stored"
            }

            telegram.ReplyToMessage(chatId, messageId, responseText)
          }()
        case "get_cat":
          go func() {
            chatId, messageId, catPhoto, err := botcommands2.GetCat(update, parsedMessage[2])

            if err != nil || len(catPhoto) == 0 {
              telegram.ReplyToMessage(chatId, messageId, "Wrong cat")
            } else {
              telegram.SendPhoto(chatId, messageId, fmt.Sprintf("%s%s", catPhoto, "?0"))
            }
          }()
        case "cats":
          go func() {
            chatIdToReply, cats := botcommands2.GetCats(update)

            if cats != nil {
              telegram.ReplyToMessage(chatIdToReply, update.MessageId, fmt.Sprintf("Available cats are: %s", strings.Join(cats, ",")))
            }
          }()
        default:
          log.Println("unknown command: ", parsedMessage[1])
        }
      }
    }
  }
}
