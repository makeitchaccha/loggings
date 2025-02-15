package main

import (
	"fmt"
	"os"

	"github.com/disgoorg/disgo"
	"github.com/makeitchaccha/loggings/internal/pkg/command"
)

func main() {

	client, err := disgo.New(os.Getenv("DISCORD_TOKEN"))

	if err != nil {
		panic(err)
	}

	commands := []command.Command{
		&command.SettingsCommand{},
	}

	for _, cmd := range commands {
		_, err := client.Rest().CreateGlobalCommand(client.ApplicationID(), cmd.Create())
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("successfully deployed commands")
}
