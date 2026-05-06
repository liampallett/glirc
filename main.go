package main

import (
	"log"
	"net"

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

		var msg Message
		msg, err = client.parseInput(text)
		if err != nil {
			client.print("%s\n", err)
			return
		}
		if msg.command != "" {
			err = client.send(msg)
			if msg.command == "QUIT" {
				client.ui.App.Stop()
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			if msg.command == "PRIVMSG" {
				echo := Message{client.nick, msg.command, msg.parameters}
				client.handlePrivmsg(echo)
			}
		}
	})

	if err = client.connect(); err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client.conn)
	if err = client.register(); err != nil {
		log.Fatal(err)
	}

	go client.readLoop()
	if err = ui.App.Run(); err != nil {
		log.Fatal(err)
	}
}
