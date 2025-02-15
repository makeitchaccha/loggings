package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/makeitchaccha/loggings/internal/app/bot"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("settings.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	database = db
}

func main() {
	bot, err := bot.New(os.Getenv("DISCORD_TOKEN"), database)

	if err != nil {
		log.Fatal("error while creating disgo client: ", err)
	}

	if err := bot.Open(context.TODO()); err != nil {
		log.Fatal("error while opening gateway: ", err)
	}

	defer bot.Close(context.TODO())

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
