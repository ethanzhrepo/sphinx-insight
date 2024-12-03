package main

import (
	"log"
	"os"

	"github.com/ethanzhrepo/sphinx-insight/core/db"
	"github.com/ethanzhrepo/sphinx-insight/core/notifier"
	pipeline "github.com/ethanzhrepo/sphinx-insight/core/pipline"
	"github.com/ethanzhrepo/sphinx-insight/core/processor"
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

	telegram := notifier.NewTelegramBot()
	go telegram.Connect(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
	)

	processors := []processor.Processor{
		// processor.NewChatgptProcessor(),
		processor.NewTelegramProcessor(telegram),
	}
	pip := pipeline.NewPipeline(processors...)

	pubsub := notifier.NewSimplePubSub()
	ch := pubsub.Subscribe(task.BinanceAnnouncement)

	go pip.Start(ch)

	binanceTask := task.NewBinanceTask(pubsub, db)
	defer binanceTask.Close()

	binanceTask.Run()
}
