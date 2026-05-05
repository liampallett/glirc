package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	ui := initUI()
	client := NewClient(config.Nick, config.User, config.Server, config.Port, ui)

	client.ui.Input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}

		text := client.ui.Input.GetText()
		client.ui.Input.SetText("")

		msg, err := client.parseInput(text)
		if err != nil {
			client.print("%s\n", err)
			return
		}
		if msg.command != "" {
			client.send(msg)
			if msg.command == "PRIVMSG" {
				echo := Message{client.nick, msg.command, msg.parameters}
				client.handlePrivmsg(echo)
			}
		}
	})

	if err = client.connect(); err != nil {
		log.Fatal(err)
	}
	defer client.conn.Close()
	client.register()

	go client.readLoop()
	if err = ui.App.Run(); err != nil {
		log.Fatal(err)
	}
}
