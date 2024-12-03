package notifier

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// TelegramBot

type TelegramBot struct {
	bot *tgbotapi.BotAPI
}

// NewTelegramBot

func NewTelegramBot() *TelegramBot {
	return &TelegramBot{}
}

// Connect to telegram
func (t *TelegramBot) Connect(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	t.bot = bot

	if err != nil {
		return err
	}

	if isDebug() {
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if isDebug() {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		}
	}

	return nil
}

// SendText
func (tp *TelegramBot) SendText(chatID int64, text string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "html"
	if isDebug() {
		log.Println("Sending message: ", text)
	}
	// if msg is empty return
	if text == "" {
		log.Println("Message is empty")
		return tgbotapi.Message{}, nil
	}
	return tp.bot.Send(msg)
}

// isDebug
func isDebug() bool {
	return os.Getenv("DEBUG") == "true"
}
