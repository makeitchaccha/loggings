package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/makeitchaccha/loggings/internal/app/bot"
	"github.com/makeitchaccha/loggings/internal/pkg/config"
	"gorm.io/gorm"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "config.yml", "config file")
	flag.Parse()
}

func main() {

	config := config.New(configFile)

	db, err := gorm.Open(config.DatabaseDialector(), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	bot, err := bot.New(os.Getenv("DISCORD_TOKEN"), db)

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
