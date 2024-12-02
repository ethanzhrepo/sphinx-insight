package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/ethanzhrepo/sphinx-insight/core/db"
	"github.com/ethanzhrepo/sphinx-insight/core/notifier"
	"github.com/ethanzhrepo/sphinx-insight/core/task"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := db.OpenDB("data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// telegram
	// from env
	bot_token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chat_id_str := os.Getenv("TELEGRAM_CHAT_ID")
	chat_id, err := strconv.ParseInt(chat_id_str, 10, 64)
	if err != nil {
		panic(err)
	}

	telegram_notifier := notifier.NewTelegramNotifier()
	go telegram_notifier.Connect(bot_token)

	// start binance job
	ps := notifier.NewSimplePubSub()

	telegramChannel := ps.Subscribe(task.BinanceAnnouncement)
	go func() {
		for msg := range telegramChannel {
			jsonMap := make(map[string]string)
			err := json.Unmarshal([]byte(msg), &jsonMap)
			if err != nil {
				log.Println(err)
				continue
			}
			telegram_notifier.SendTextWithLink(chat_id, jsonMap["title"], jsonMap["link"])
		}
	}()

	binanceTask := task.NewBinanceTask(ps, db)
	binanceTask.Run()
}
