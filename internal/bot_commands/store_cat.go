package bot_commands

import (
	"awesomeProject/internal/bot"
	"errors"
	"net/http"
	"regexp"
)

var availableExtensions = map[string]bool{"jpeg": true, "png": true}

func ExtractPhotoId(message *bot.Message) string {

	if message != nil {
		photosAmount := len(message.Photo)

		if photosAmount != 0 {
			return message.Photo[photosAmount-1].FileId
		}
	}

	return ""
}

func ValidateCatFile(fileExtension string) bool {
	val, ok := availableExtensions[fileExtension]
	return ok && val
}

func ExtractFileExtension(file []byte) string {
	contentType := http.DetectContentType(file)
	rg := regexp.MustCompile("image/(jpg|jpeg|png)")
	match := rg.FindStringSubmatch(contentType)

	if len(match) < 2 {
		return ""
	}

	return match[1]
}

// StoreCat TODO add chats whitelist for uploading images
func StoreCat(cat string, image []byte) error {
	catExtension := ExtractFileExtension(image)

	if len(catExtension) == 0 {
		return errors.New("wrong file extension")
	}

	isFileValid := ValidateCatFile(catExtension)

	if !isFileValid {
		return errors.New("wrong file extension")
	}

	return GetService().StoreCat(cat, catExtension, image)
}
