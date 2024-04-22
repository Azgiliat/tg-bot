package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Payload map[string]interface{}

type File struct {
	FileId   string `json:"file_id"`
	FilePath string `json:"file_path"`
}

type Photo struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

type From struct {
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Chat struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type Entity struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}

type Message struct {
	MessageId       int    `json:"message_id"`
	Date            int    `json:"date"`
	Text            string `json:"text"`
	From            `json:"from"`
	Chat            `json:"chat"`
	Photo           []Photo  `json:"photo"`
	Entities        []Entity `json:"entities"`
	CaptionEntities []Entity `json:"caption_entities"`
	Caption         string   `json:"caption"`
	ReplyToMessage  *Message `json:"reply_to_message"`
}

type Update struct {
	UpdateId    int `json:"update_id"`
	Message     `json:"message"`
	Description string `json:"description"`
}

type FileResponse struct {
	Ok     bool `json:"ok"`
	Result File `json:"result"`
}

type UpdatesResponse struct {
	Ok      bool     `json:"ok"`
	Results []Update `json:"result"`
}

type TgBot struct {
	token                string
	apiUrl               string
	fileApiUrl           string
	lastUpdateIdentifier int
}

func NewBot(token string) *TgBot {
	return &TgBot{token: token, fileApiUrl: fmt.Sprintf("https://api.telegram.org/file/bot%s", token), apiUrl: fmt.Sprintf("https://api.telegram.org/bot%s", token)}
}

func (bot *TgBot) sendPostRequest(url string, payload Payload) ([]byte, error) {
	body, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(body)
	res, err := http.Post(url, "application/json", buff)
	defer res.Body.Close()

	if err != nil {
		log.Println("post request failed")

		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Println(res.Status)

		return nil, err
	}

	body, err = io.ReadAll(res.Body)

	if err != nil {
		log.Println("Failed to read response body")

		return nil, err
	}

	return body, nil
}

func (bot *TgBot) sendGetRequest(url string) ([]byte, error) {
	res, err := http.Get(url)
	statusCode := res.StatusCode
	defer res.Body.Close()

	if err != nil || statusCode != http.StatusOK {
		log.Println("get request failed")

		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println("Failed to read response body")

		return nil, err
	}

	return body, nil
}

func (bot *TgBot) GetUpdates() []Update {
	url := fmt.Sprintf("%s/getUpdates?timeout=%d&offset=%d", bot.apiUrl, time.Minute/time.Millisecond, bot.lastUpdateIdentifier)
	response, err := bot.sendGetRequest(url)

	if err != nil {
		log.Println("failed to get updates")

		return nil
	}

	result := UpdatesResponse{}
	err = json.Unmarshal(response, &result)

	if err != nil || !result.Ok {
		return nil
	}

	return result.Results
}

func (bot *TgBot) callSendMessageEndpoint(payload Payload) ([]byte, error) {
	url := fmt.Sprintf("%s/sendMessage", bot.apiUrl)
	res, err := bot.sendPostRequest(url, payload)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (bot *TgBot) SendMessage(chatId int, message string) {
	payload := map[string]interface{}{
		"chat_id": chatId,
		"text":    message,
	}

	_, err := bot.callSendMessageEndpoint(payload)

	if err != nil {
		log.Println("failed to send message")
	}
}

func (bot *TgBot) ReplyToMessage(chatId int, messageId int, message string) {
	payload := Payload{
		"chat_id": chatId,
		"text":    message,
		"reply_parameters": Payload{
			"message_id": messageId,
		},
	}

	_, err := bot.callSendMessageEndpoint(payload)

	if err != nil {
		log.Println("failed to send message")
	}
}

func (bot *TgBot) SendPhoto(chatId int, messageId int, photoId string) {
	payload := Payload{
		"chat_id": chatId,
		"photo":   photoId,
	}

	if messageId != 0 {
		payload["reply_parameters"] = Payload{
			"message_id": messageId,
		}
	}

	_, err := bot.sendPostRequest(fmt.Sprintf("%s/sendPhoto", bot.apiUrl), payload)

	if err != nil {
		log.Println("failed to send photo")
	}
}

func (bot *TgBot) GetFile(fileId string) (*File, error) {
	var file = FileResponse{}
	res, err := bot.sendGetRequest(fmt.Sprintf("%s/getFile?file_id=%s", bot.apiUrl, fileId))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &file)

	if err != nil {
		return nil, err
	}

	return &file.Result, nil
}

func (bot *TgBot) DownloadFile(fileId string) ([]byte, error) {
	file, err := bot.GetFile(fileId)

	if err != nil || len(file.FilePath) == 0 {
		return nil, err
	}

	res, err := bot.sendGetRequest(fmt.Sprintf("%s/%s", bot.fileApiUrl, file.FilePath))

	return res, nil
}

func (bot *TgBot) GetUpdatesChannel() <-chan Update {
	channel := make(chan Update, 100)

	go func() {
		for {
			updates := bot.GetUpdates()

			if updates == nil {
				log.Println("Sleep 5 sec before next updates")
				time.Sleep(5 * time.Second)

				continue
			}

			for _, update := range updates {
				bot.lastUpdateIdentifier = update.UpdateId + 1
				channel <- update
			}
		}
	}()

	return channel
}
