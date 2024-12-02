package notifier

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// telegramNotifier

type TelegramNotifier struct {
	bot *tgbotapi.BotAPI
}

// NewTelegramNotifier

func NewTelegramNotifier() *TelegramNotifier {
	return &TelegramNotifier{}
}

// Connect to telegram
func (t *TelegramNotifier) Connect(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	t.bot = bot

	if err != nil {
		return err
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

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

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}

	return nil
}

// SendText
func (t *TelegramNotifier) SendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	t.bot.Send(msg)
}

// SendText With Link
func (t *TelegramNotifier) SendTextWithLink(chatID int64, text string, link string) {
	msg := tgbotapi.NewMessage(chatID, text)
	// html
	msg.ParseMode = "html"
	msg.Text = "<a href=\"" + link + "\">" + text + "</a>"
	t.bot.Send(msg)
}
