package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	ui := initUI()
	client := NewClient(config.Nick, config.User, config.Server, config.Port, ui)

	client.ui.Channels.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		client.currentChannel = mainText
		client.ui.Chat.SetTitle(mainText)
		client.refreshNames()
		client.ui.Chat.Clear()
		ch, ok := client.channels[mainText]
		if !ok {
			return
		}
		_, err = fmt.Fprint(client.ui.Chat, strings.Join(ch.history, ""))
		if err != nil {
			return
		}
		client.ui.App.SetFocus(client.ui.Input)
	})

	client.ui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if client.ui.App.GetFocus() == client.ui.Input {
				client.ui.App.SetFocus(client.ui.Channels)
			} else {
				client.ui.App.SetFocus(client.ui.Input)
			}
			return nil
		}
		return event
	})
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
			err = client.server.send(msg)
			if msg.command == "QUIT" {
				client.ui.App.Stop()
				return
			}
			if err != nil {
				return
			}
			if msg.command == "PRIVMSG" {
				echo := Message{client.nick, msg.command, msg.parameters}
				client.handlePrivmsg(echo)
			}
		}
	})

	if err = client.server.connect(); err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client.server.conn)
	if err = client.register(); err != nil {
		log.Fatal(err)
	}

	go client.server.readLoop()
	if err = ui.App.Run(); err != nil {
		log.Fatal(err)
	}
}
