package processor

import (
	"log"
	"os"
	"strconv"

	"github.com/ethanzhrepo/sphinx-insight/core/notifier"
)

type TelegramProcessor struct {
	chat_id  int64
	notifier *notifier.TelegramBot
}

func NewTelegramProcessor(
	notifier *notifier.TelegramBot,
) *TelegramProcessor {
	tp := TelegramProcessor{}
	chat_id_str := os.Getenv("TELEGRAM_CHAT_ID")
	chat_id, err := strconv.ParseInt(chat_id_str, 10, 64)
	if err != nil {
		panic(err)
	}
	tp.chat_id = chat_id
	tp.notifier = notifier
	return &tp
}

func (tp *TelegramProcessor) Process(data *ProcessorData) (string, error) {

	_, err := tp.notifier.SendText(tp.chat_id, data.Content)
	if err != nil {
		log.Println("Error sending text to telegram", err)
		return "", err
	}
	return data.Content, err
}

func (tp *TelegramProcessor) Name() string {
	return "TelegramProcessor"
}
